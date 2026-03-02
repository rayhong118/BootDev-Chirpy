package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rayhong118/BootDev-Chirpy/internal/auth"
	"github.com/rayhong118/BootDev-Chirpy/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type createUserResponse struct {
		Email     string    `json:"email"`
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	request := createUserRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	email := request.Email
	hashedPassword, err := auth.HashPassword(request.Password)

	user, createUserErr := cfg.db.CreateUser(r.Context(), database.CreateUserParams{email, hashedPassword})

	if createUserErr != nil {
		respondWithError(w, 500, "Create user error", err)
		return
	}

	respondWithJSON(w, 201, createUserResponse{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Id:        user.ID,
	})

}
