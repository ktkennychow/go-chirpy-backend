package main

import "net/http"

func main(){
	serveMux := http.NewServeMux()
	server := http.Server{Handler: serveMux, Addr: "localhost:8080" }
	
	dir := http.Dir(".")
	fileServer := http.FileServer(dir)

	serveMux.Handle("/app/*", http.StripPrefix("/app", fileServer))

	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	server.ListenAndServe()
}