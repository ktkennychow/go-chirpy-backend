package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main(){
	const filepathRoot = "."
	const port = "8080"
	
	apiConfig := apiConfig{fileserverHits: 0}

	sMux := http.NewServeMux()

	server := &http.Server{Handler: sMux, Addr: ":" + port }
	
	dir := http.Dir(".")
	handlerfs := apiConfig.middlewareMetricsInc(http.FileServer(dir))

	sMux.Handle("/app/*", http.StripPrefix("/app", handlerfs))

	sMux.HandleFunc("GET /healthz", handlerReadiness)

	sMux.HandleFunc("GET /metrics", apiConfig.handlerMetrics)

	sMux.HandleFunc("/reset", apiConfig.handlerReset)

	log.Printf("Serving files from %v on port: %v", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}