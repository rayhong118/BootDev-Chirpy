package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type successResponse struct {
		Cleaned string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
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

	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	respondWithJSON(w, 200, successResponse{
		Cleaned: cleanChirp(params.Body, profane),
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
