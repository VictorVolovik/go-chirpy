package main

import (
	"encoding/json"
	"net/http"
	"time"

	"VictorVolovik/go-chirpy/internal/auth"
	"VictorVolovik/go-chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	accessTokenExpirationTime := time.Hour
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, accessTokenExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to login", err)
		return
	}

	refreshTokenExpirationTimestamp := time.Now().Add(time.Hour * 24 * 60) // in 60 days
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to login", err)
		return
	}
	err = cfg.db.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refreshToken,
			ExpiresAt: refreshTokenExpirationTimestamp,
			UserID:    user.ID,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to login", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
