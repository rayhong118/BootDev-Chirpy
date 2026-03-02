package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rayhong118/BootDev-Chirpy/internal/auth"
)

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type LoginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type LoginResponse struct {
		Email     string    `json:"email"`
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := LoginPayload{}
	err := decoder.Decode(&payload)

	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	email := payload.Email
	hashedPassword, err := auth.HashPassword(payload.Password)

	user, getUserErr := cfg.db.GetUserByEmail(r.Context(), email)

	if getUserErr != nil {
		respondWithError(w, 401, "Incorrect email ", getUserErr)
		return
	}

	storedHashedPassword := user.HashedPassword

	if hashedPassword != storedHashedPassword {
		respondWithError(w, 401, "Incorrect email or password", getUserErr)
		return
	}

	respondWithJSON(w, 200, LoginResponse{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Id:        user.ID,
	})

}
