package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const MAXDURATION = 1 * time.Hour

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
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	user, err := cfg.DB.ReadSingleUserbyEmail(reqBody.Email)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 401)
		return
	}
	
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(reqBody.Password))
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 401)
		return
	}

	var tokenLife int

	if reqBody.Expires_in_seconds == 0 || reqBody.Expires_in_seconds > int(MAXDURATION.Seconds()) {
		tokenLife = int(MAXDURATION.Seconds())
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
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	random32Bytes := make([]byte, 32)
	
	_, err = rand.Read([]byte(random32Bytes))
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	refreshTokenString := hex.EncodeToString(random32Bytes)
	refreshTokenExpiry := time.Now().UTC().Add(60 * time.Hour)

	refreshToken, err := cfg.DB.CreateRefreshTokenWDetails(user.ID, refreshTokenString, refreshTokenExpiry)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	updatedUser, err := cfg.DB.UpdateUser(user.Email, user.HashedPassword, user.ID)
	if err != nil {
		cfg.DB.DeleteRefreshToken(refreshTokenString)
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	respBody.ID = updatedUser.ID
	respBody.Email = updatedUser.Email
	respBody.Token = signedJwtToken
	respBody.RefreshToken = refreshToken.RefreshToken
	respBody.IsChirpyRed = updatedUser.IsChirpyRed
	
	dat, err := json.Marshal(respBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerRefreshAuth(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	refreshToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	refreshTokenStruct, err := cfg.DB.ReadSingleRefreshTokenWDetails(refreshToken)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 401)
		return
	}

	tokenLife := int(MAXDURATION.Seconds())
	expirationTime := time.Now().UTC().Add(time.Duration(tokenLife) * time.Second)

	newJwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   fmt.Sprint(refreshTokenStruct.UserID),
	})

	signedJwtToken, err := newJwtToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}
	
	random32Bytes := make([]byte, 32)
	
	_, err = rand.Read([]byte(random32Bytes))
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	respBody.Token = signedJwtToken

	dat, err := json.Marshal(respBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerRevokeAuth(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type", "application/json")
	respBody := &RespBody{}

	refreshToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	err := cfg.DB.DeleteRefreshToken(refreshToken)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 401)
		return
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		cfg.handlerErrors(w, err, respBody, 500)
		w.Write(dat)
		return
	}

	w.WriteHeader(204)
}

func (cfg *apiConfig)handlerAuthenticateWJwt(r *http.Request)(int, error){
	var userID int

	jwtTokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	jwtToken, err := jwt.ParseWithClaims(jwtTokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {return []byte(cfg.jwtSecret), nil})

	if err != nil {
		return userID, err
	}

	idString, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return userID, err
	}
	userID, err = strconv.Atoi(idString)
	if err != nil {
		return userID, err
	}
	return userID, nil
}