package main

import (
	"encoding/json"
	"net/http"

	"VictorVolovik/go-chirpy/internal/auth"
	"VictorVolovik/go-chirpy/internal/database"
)

func (cfg *apiConfig) handleUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
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

	if len(params.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "Email field should not be empty", nil)
		return
	}

	if len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "Password field should not be empty", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user data", err)
		return
	}

	user, err := cfg.db.UpdateUserEmailAndPassword(
		r.Context(),
		database.UpdateUserEmailAndPasswordParams{
			ID:             userID,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user data", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
