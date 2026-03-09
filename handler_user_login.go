package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rayhong118/BootDev-Chirpy/internal/auth"
	"github.com/rayhong118/BootDev-Chirpy/internal/database"
)

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type LoginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type LoginResponse struct {
		Email        string    `json:"email"`
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := LoginPayload{}
	err := decoder.Decode(&payload)

	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	email := payload.Email

	user, getUserErr := cfg.db.GetUserByEmail(r.Context(), email)

	if getUserErr != nil {
		respondWithError(w, 401, "Incorrect email or password", getUserErr)
		return
	}

	match, checkPasswordErr := auth.CheckPasswordHash(payload.Password, user.HashedPassword)

	if checkPasswordErr != nil {
		respondWithError(w, 500, "Something went wrong", checkPasswordErr)
		return
	}

	if match != true {
		respondWithError(w, 401, "Incorrect email or password", errors.New("Incorrect email or password"))
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.Secret)

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})

	if err != nil {
		respondWithError(w, 500, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, 200, LoginResponse{
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Id:           user.ID,
		Token:        token,
		RefreshToken: refreshToken,
	})

}
