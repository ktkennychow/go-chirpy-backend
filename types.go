package main

import (
	"sync"
	"time"
)

type RespBody struct {
	ID int `json:"id"`
	Error string `json:"error"`
	Body string `json:"body"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	HashedPassword []byte `json:"hashed_password"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
	RefreshTokens map[string]RefreshToken
}

type RefreshToken struct {
	UserID int
	RefreshToken string
	ExpiresAt time.Time
}
