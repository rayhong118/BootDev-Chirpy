package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
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
	Email       string    `json:"email"`
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// handle user email and password update
func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	token, getTokenErr := auth.GetBearerToken(r.Header)

	if getTokenErr != nil {
		respondWithError(w, 401, "Failed to get token from header", getTokenErr)
		return

	}
	userId, tokenErr := auth.ValidateJWT(token, cfg.Secret)

	if tokenErr != nil {
		respondWithError(w, 401, "Failed to validate token from header", tokenErr)
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
		Email:       newUserInfo.Email,
		CreatedAt:   newUserInfo.CreatedAt,
		UpdatedAt:   newUserInfo.UpdatedAt,
		Id:          newUserInfo.ID,
		IsChirpyRed: newUserInfo.IsChirpyRed.Bool,
	})

}

type subscriptionUpdatePayload struct {
	Event string                 `json:"event"`
	Data  subscriptionUpdateData `json:"data"`
}

type subscriptionUpdateData struct {
	UserID uuid.UUID `json:"user_id"`
}

// handle subscription update. it only add subscription for now
func (cfg *apiConfig) handleSubscriptionUpdate(w http.ResponseWriter, r *http.Request) {
	apiKey, getApiKeyErr := auth.GetPolkaKey(r.Header)

	if getApiKeyErr != nil || apiKey != os.Getenv("POLKA_KEY") {
		respondWithError(w, 401, "Incorrect or missing API key", getApiKeyErr)
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := subscriptionUpdatePayload{}
	err := decoder.Decode(&request)

	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	event := request.Event

	if event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	updateErr := cfg.db.UpdateUserRedSubscription(r.Context(), database.UpdateUserRedSubscriptionParams{ID: request.Data.UserID, IsChirpyRed: sql.NullBool{Bool: true, Valid: true}})

	if updateErr != nil {
		respondWithError(w, 404, "User can't be found", updateErr)
		return
	}

	respondWithJSON(w, 204, nil)

}
