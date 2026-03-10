package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rayhong118/BootDev-Chirpy/internal/auth"
	"github.com/rayhong118/BootDev-Chirpy/internal/database"
)

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserResponse struct {
	Email     string    `json:"email"`
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	token, getTokenErr := auth.GetBearerToken(r.Header)

	if getTokenErr != nil {
		respondWithError(w, 401, "Failed to get token from header", getTokenErr)
		return

	}
	userId, tokenErr := auth.ValidateJWT(token, cfg.Secret)

	if tokenErr != nil {
		respondWithError(w, 401, "Failed to validate token from header", getTokenErr)
		return
	}

	decoder := json.NewDecoder(r.Body)
	updatePayload := updateUserRequest{}
	err := decoder.Decode(&updatePayload)

	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	hashedPassword, hashPwdErr := auth.HashPassword(updatePayload.Password)

	if hashPwdErr != nil {
		respondWithError(w, 500, "Failed to hash password", hashPwdErr)
		return
	}

	newUserInfo, updateUserErr := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          updatePayload.Email,
		HashedPassword: hashedPassword,
	})

	if updateUserErr != nil {
		respondWithError(w, 500, "Failed to update user", updateUserErr)
		return
	}

	respondWithJSON(w, 200, updateUserResponse{
		Email:     newUserInfo.Email,
		CreatedAt: newUserInfo.CreatedAt,
		UpdatedAt: newUserInfo.UpdatedAt,
		Id:        newUserInfo.ID,
	})

}
