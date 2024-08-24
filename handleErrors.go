package main

import (
	"fmt"
	"net/http"
)

func handlerErrors(w http.ResponseWriter, err error, respBody *RespBody, code int) {
	fmt.Println(code, err)
	w.WriteHeader(code)
	respBody.Error = err.Error()
}
