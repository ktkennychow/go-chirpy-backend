package main

import (
	"net/http"
)

func handlerErrors(w http.ResponseWriter, err error, respBody *RespBody, code int) {
	w.WriteHeader(code)
	respBody.Error = err.Error()
}
