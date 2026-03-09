package main

import (
	"net/http"
	"time"

	"github.com/rayhong118/BootDev-Chirpy/internal/auth"
)

type TokenRefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	refreshToken, getTokenErr := auth.GetBearerToken(header)

	if getTokenErr != nil {
		respondWithError(w, 401, "Failed to get token from header", getTokenErr)
		return
	}

	authToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)

	if err != nil {
		respondWithError(w, 401, "No Token", err)
		return
	}

	if authToken.RevokedAt.Valid {
		respondWithError(w, 401, "Revoked", nil)
		return
	}

	if authToken.ExpiresAt.Valid && authToken.ExpiresAt.Time.Before(time.Now()) {
		respondWithError(w, 401, "Expired", nil)
		return
	}

	newToken, newTokenErr := auth.MakeJWT(authToken.UserID.UUID, cfg.Secret)

	if newTokenErr != nil {
		respondWithError(w, 401, "Failed to create new token", nil)
		return
	}

	respondWithJSON(w, 200, TokenRefreshResponse{Token: newToken})

}

func (cfg *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {

	header := r.Header
	refreshToken, getTokenErr := auth.GetBearerToken(header)
	if getTokenErr != nil {
		respondWithError(w, 401, "Failed to get token from header", getTokenErr)
	}

	cfg.db.RevokeRefreshToken(r.Context(), refreshToken)

	respondWithJSON(w, 204, nil)
}
