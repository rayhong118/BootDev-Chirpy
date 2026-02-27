package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	apiCfg := apiConfig{}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func validateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	type successResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder((r.Body))
	params := parameters{}
	err := decoder.Decode(&params)

	defaultErrBody := errorResponse{
		Error: "Something went wrong",
	}
	defaultErr, defaultErrEncodeErr := json.Marshal(defaultErrBody)

	if err != nil || defaultErrEncodeErr != nil {
		w.WriteHeader(500)
		w.Write(defaultErr)
		return
	}
	if len(params.Body) > 140 {

		responseBody := errorResponse{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			w.WriteHeader(500)
			w.Write(defaultErr)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	} else {
		responseBody := successResponse{
			Valid: true,
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			w.WriteHeader(500)
			w.Write(defaultErr)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
		return
	}

}
