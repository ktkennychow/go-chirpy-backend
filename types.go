package main

type RespBody struct {
	ID int `json:"id"`
	Error string `json:"error"`
	Body string `json:"body"`
	Email string `json:"email"`
	Token string `json:"token"`
}