package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	outputHTML(w, "metrics/index.html", cfg)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	cfg.FileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type returnVals struct {
		Error string `json:"error"`
		Valid bool `json:"valid"`
		Cleaned_body string `json:"cleaned_body"`
	}
	respBody := &returnVals{}

	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		respBody.Error = "Something went wrong"
		return
	} 
	
	if len(params.Body) > 140 {
		respBody.Error = "Chirp is too long"
		w.WriteHeader(400)
		} else {
		bannedWords := []string{"kerfuffle","sharbert","fornax"}
		words := strings.Split(params.Body, " ")
		for i, word := range words {
			if slices.Contains(bannedWords, strings.ToLower(word)) {
				words[i] = "****"
			}
		}
		respBody.Cleaned_body = strings.Join(words, " ")
		w.WriteHeader(200)
	}
	

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		respBody.Error = "Something went wrong"
		return
	}
	
	w.Write(dat)
}