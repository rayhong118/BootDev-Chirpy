package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
		UserID: chirp.UserId,
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
