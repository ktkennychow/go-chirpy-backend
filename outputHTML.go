package main

import (
	"html/template"
	"net/http"
)

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
    t, err := template.ParseFiles(filename)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    if err := t.Execute(w, data); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}