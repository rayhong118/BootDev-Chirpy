package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rayhong118/BootDev-Chirpy/internal/auth"
	"github.com/rayhong118/BootDev-Chirpy/internal/database"
)

type PostChirpBody struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	type successResponse struct {
		Cleaned string `json:"cleaned_body"`
	}

	tokenString, tokenErr := auth.GetBearerToken(r.Header)

	if tokenErr != nil {
		respondWithError(w, 401, "401 Unauthorized tokenErr", tokenErr)
		return
	}

	userID, validationErr := auth.ValidateJWT(tokenString, cfg.Secret)

	if validationErr != nil {
		respondWithError(w, 401, "401 Unauthorized validationErr", validationErr)
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := PostChirpBody{}
	err := decoder.Decode(&chirp)

	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {

		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	responseChirp, err := cfg.db.SaveChirp(r.Context(), database.SaveChirpParams{
		Body:   cleanChirp(chirp.Body, profane),
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 500, "Could not create chirp", err)
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        responseChirp.ID,
		CreatedAt: responseChirp.CreatedAt,
		UpdatedAt: responseChirp.UpdatedAt,
		Body:      responseChirp.Body,
		UserId:    responseChirp.UserID,
	})
}

func cleanChirp(chirp string, profane map[string]struct{}) string {

	chirpSlice := strings.Fields(chirp)
	for index, word := range chirpSlice {
		if _, ok := profane[strings.ToLower(word)]; ok {
			chirpSlice[index] = "****"
		}
	}
	return strings.Join(chirpSlice, " ")
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	authorId := r.URL.Query().Get("author_id")
	authorUUID := uuid.NullUUID{}
	if authorId != "" {
		u, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, 400, "Invalid author_id", err)
			return
		}
		authorUUID.UUID = u
		authorUUID.Valid = true
	}

	chirps, getChirpsErr := cfg.db.GetChirps(r.Context(), authorUUID)

	if getChirpsErr != nil {
		respondWithError(w, 500, "Chirp fetch failed", getChirpsErr)
		return
	}
	output := make([]Chirp, len(chirps))

	for i, chirp := range chirps {
		output[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
	}

	respondWithJSON(w, 200, output)

}

func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")

	chirpUUID, parseErr := uuid.Parse(chirpId)
	if parseErr != nil {
		respondWithError(w, 404, "Chirp fetch failed", parseErr)
		return
	}

	chirp, getChirpErr := cfg.db.GetChirpByID(r.Context(), chirpUUID)

	if getChirpErr != nil {
		respondWithError(w, 404, "Chirp fetch failed", getChirpErr)
		return
	}

	respondWithJSON(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) handleDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")

	token, getTokenError := auth.GetBearerToken(r.Header)
	if getTokenError != nil {
		respondWithError(w, 401, "unauthed", getTokenError)
		return
	}

	userId, tokenErr := auth.ValidateJWT(token, cfg.Secret)
	if tokenErr != nil {
		respondWithError(w, 401, "Failed to validate token from header", tokenErr)
		return
	}

	chirpUUID, parseErr := uuid.Parse(chirpId)
	if parseErr != nil {
		respondWithError(w, 403, "Invalid chirp ID", parseErr)
		return
	}

	chirp, getChirpErr := cfg.db.GetChirpByID(r.Context(), chirpUUID)

	if getChirpErr != nil {
		respondWithError(w, 404, "Cannot get chirp", getChirpErr)
		return
	}
	if chirp.UserID != userId {
		respondWithError(w, 403, "Not the correct user", errors.New("Not the correct user"))
		return
	}

	deleteChirpErr := cfg.db.DeleteChirpById(r.Context(), chirp.ID)

	if deleteChirpErr != nil {
		respondWithError(w, 500, "Failed to delete chirp by ID", deleteChirpErr)
		return
	}

	respondWithJSON(w, 204, nil)

}
