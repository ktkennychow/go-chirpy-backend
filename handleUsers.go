package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	type reqParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 1)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	user, err := cfg.DB.CreateUsers(reqBody.Email, hashedPassword)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	respBody.ID = user.ID
	respBody.Email = user.Email
	respBody.IsChirpyRed = false
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerModifyUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	userID, err := cfg.handlerAuthenticateWJwt(r)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 401)
		return
	}

	type reqParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 1)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	user, err := cfg.DB.UpdateUser(reqBody.Email, hashedPassword, userID)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	respBody.ID = user.ID
	respBody.Email = user.Email
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerPolkaUserUpgrade(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	type reqParams struct {
		Event string `json:"event"`
		Data struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	if reqBody.Event == "user.upgraded" {
		_, err := cfg.DB.UpgradeUser(reqBody.Data.UserID)
		if err != nil {
			if err.Error() == "User does not exist"{
				cfg.handlerErrors(w, err, respBody, 404)
			}
			cfg.handlerErrors(w, err, respBody, 500)
			return
		}
	}
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(204)
	w.Write(dat)
}