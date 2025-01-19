package main

import (
	"VictorVolovik/go-chirpy/internal/database"
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsGetById(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", nil)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	authorID := query.Get("author_id")

	var chirps []database.Chirp
	var err error
	if authorID == "" {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting all chirps", err)
			return
		}
	} else {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByUserId(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps for user %s", authorID), err)
			return
		}
	}

	returnChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		returnChirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
	}

	if sortOrder := query.Get("sort"); sortOrder == "desc" {
		sort.Slice(returnChirps, func(i, j int) bool {
			return returnChirps[i].CreatedAt.After(returnChirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, returnChirps)
}
