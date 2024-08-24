package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	type reqParams struct {
		Password string `json:"password"`
		Email string `json:"email"`
		Expires_in_seconds int `json:"expires_in_seconds"`
	}
	reqBody := reqParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	user, err := cfg.DB.ReadSingleUserbyEmail(reqBody.Email)
	if err != nil {
		handlerErrors(w, err, respBody, 401)
		return
	}
	
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(reqBody.Password))
	if err != nil {
		handlerErrors(w, err, respBody, 401)
		return
	}

	var tokenLife int
	const duration = 24 * time.Hour

	if reqBody.Expires_in_seconds == 0 || reqBody.Expires_in_seconds > int(duration.Seconds()) {
		tokenLife = int(duration.Seconds())
	} else {
		tokenLife = reqBody.Expires_in_seconds
	}

	expirationTime := time.Now().UTC().Add(time.Duration(tokenLife) * time.Second)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   fmt.Sprint(user.ID),
	})

	signedJwtToken, err := jwtToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	respBody.ID = user.ID
	respBody.Email = user.Email
	respBody.Token = signedJwtToken
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		handlerErrors(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(200)
	w.Write(dat)
}