package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"VictorVolovik/go-chirpy/internal/auth"
	"VictorVolovik/go-chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No JWT provided", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	replacer := "****"
	cleaned := replaceWords(params.Body, badWords, replacer)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleaned, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new chirp", nil)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func replaceWords(str string, wordsToReplace map[string]struct{}, replacer string) string {
	strWords := strings.Fields(str)

	for i, word := range strWords {
		wordLowercased := strings.ToLower(word)
		if _, ok := wordsToReplace[wordLowercased]; ok {
			strWords[i] = replacer
		}
	}

	joined := strings.Join(strWords, " ")
	return joined
}
