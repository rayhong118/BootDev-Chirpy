package main

import (
	"encoding/json"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type successResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder((r.Body))
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {

		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, 200, successResponse{
		Valid: true,
	})

}
