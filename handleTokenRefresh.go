package main

import (
	"net/http"
	"time"

	"VictorVolovik/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided", err)
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	accessTokenExpirationTime := time.Hour
	accessToken, err := auth.MakeJWT(userID, cfg.jwtSecret, accessTokenExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
