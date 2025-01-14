package main

import (
	"VictorVolovik/go-chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("unable to open database connection: %s", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handleHealthCheck)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetricsCheck)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleMetricsReset)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}
