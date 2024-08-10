package main

import (
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func main(){
	serveMux := http.NewServeMux()
	server := http.Server{Handler: serveMux, Addr: "localhost:8080" }
	apiConfig := apiConfig{}
	
	dir := http.Dir(".")
	fileServer := http.FileServer(dir)

	serveMux.Handle("/app/*", http.StripPrefix("/app", apiConfig.middlewareMetricsInc(fileServer)))

	serveMux.HandleFunc("/healthz", handlerReadiness)

	serveMux.HandleFunc("/metrics", apiConfig.handlerMetrics)

	serveMux.HandleFunc("/reset", apiConfig.handlerReset)

	server.ListenAndServe()
}