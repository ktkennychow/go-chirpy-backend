package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	authorID, err := cfg.handlerAuthenticateWJwt(r)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 401)
		return
	}

	type reqParams struct {
		Body string `json:"body"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
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
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	respBody.ID = chirp.ID
	respBody.AuthorID = chirp.AuthorID
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
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
			cfg.handlerErrors(w, err, respBody, 500)
			return
		}

	authorID := r.URL.Query().Get("author_id")

	if authorID != "" {
		authorID, err := strconv.Atoi(authorID)
		if err != nil {
			cfg.handlerErrors(w, err, respBody, 500)
			return
		}
		for i, chirp := range chirps {
			if chirp.AuthorID != authorID {
				chirps = append(chirps[:i], chirps[i+1:]...)
			}	
		}
	}

	sortCommand := r.URL.Query().Get("sort")

	if sortCommand == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID 
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID 
		})
	}

	dat, err := json.Marshal(chirps)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
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
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	chirp, err := cfg.DB.ReadSingleChirp(chirpID)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 404)
		return
	}
	
	dat, err := json.Marshal(chirp)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig)handlerDeleteSingleChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	userID, err := cfg.handlerAuthenticateWJwt(r)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 401)
		return
	}
	
	chirpIDPath := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDPath)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	chirp, err := cfg.DB.ReadSingleChirp(chirpID)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 404)
		return
	}
	
	if chirp.AuthorID != userID {
		cfg.handlerErrors(w, errors.New("not authorized to delete this chirp"), respBody, 403)
		return
	}
	
	err = cfg.DB.DeleteSingleChirp(chirp.ID)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 404)
		return
	}
	
	w.WriteHeader(204)
}