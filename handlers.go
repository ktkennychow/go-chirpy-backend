package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type RespBody struct {
	Id int `json:"id"`
	Error string `json:"error"`
	Body string `json:"body"`
	Email string `json:"email"`
}

func handlerError(w http.ResponseWriter, err error, respBody *RespBody, code int) {
	log.Printf("Error: %s", err)
	w.WriteHeader(code)
	respBody.Error = "Something went wrong"
}

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

func (db *DB) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	type reqParams struct {
		Body string `json:"body"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	} 
	
	if len(reqBody.Body) > 140 {
		respBody.Error = "Chirp is too long"
		w.WriteHeader(400)
		return 
	}
	
	bannedWords := []string{"kerfuffle","sharbert","fornax"}
	words := strings.Split(reqBody.Body, " ")
	for i, word := range words {
		if slices.Contains(bannedWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	respBody.Body = strings.Join(words, " ")
	chirp, err := db.CreateChirp(respBody.Body)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}
	
	respBody.Id = chirp.ID
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(201)
	w.Write(dat)
}

func (db *DB)handlerGetChirps(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	chirps, err := db.GetChirps()
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}
	
	dat, err := json.Marshal(chirps)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}

func (db *DB)handlerGetSingleChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}
	
	chirpIDPath := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDPath)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}
	
	chirp, err := db.GetSingleChirp(chirpID)
	if err != nil {
		handlerError(w, err, respBody, 404)
		return
	}
	
	dat, err := json.Marshal(chirp)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}

func (db *DB) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	type reqParams struct {
		Email string `json:"email"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}
	
	user, err := db.CreateUsers(reqBody.Email)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}
	
	respBody.Id = user.ID
	respBody.Email = user.Email
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		handlerError(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(201)
	w.Write(dat)
}