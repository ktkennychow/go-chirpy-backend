package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	authorID, err := cfg.handlerAuthenticateWJwt(r)
	if err != nil {
		handlerErrors(w, err, respBody, 401)
		return
	}

	type reqParams struct {
		Body string `json:"body"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
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
	chirp, err := cfg.DB.CreateChirp(respBody.Body, authorID)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	respBody.ID = chirp.ID
	respBody.AuthorID = chirp.AuthorID
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig)handlerReadChirps(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	chirps, err := cfg.DB.ReadChirps()
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	dat, err := json.Marshal(chirps)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig)handlerReadSingleChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}
	
	chirpIDPath := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDPath)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	chirp, err := cfg.DB.ReadSingleChirp(chirpID)
	if err != nil {
		handlerErrors(w, err, respBody, 404)
		return
	}
	
	dat, err := json.Marshal(chirp)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig)handlerDeleteSingleChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}
	
	chirpIDPath := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDPath)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	chirp, err := cfg.DB.ReadSingleChirp(chirpID)
	if err != nil {
		handlerErrors(w, err, respBody, 404)
		return
	}
	
	dat, err := json.Marshal(chirp)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(200)
	w.Write(dat)
}