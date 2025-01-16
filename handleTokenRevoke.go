package main

import (
	"net/http"

	"VictorVolovik/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handleTokenRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
