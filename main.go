package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	FileserverHits int
}

func main(){
	const filepathRoot = "."
	const port = "8080"

	db, err := NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}
	
	apiConfig := apiConfig{FileserverHits: 0}

	sMux := http.NewServeMux()

	server := &http.Server{Handler: sMux, Addr: ":" + port }
	
	dir := http.Dir(".")
	handlerfs := apiConfig.middlewareMetricsInc(http.FileServer(dir))

	sMux.Handle("GET /app/*", http.StripPrefix("/app", handlerfs))

	sMux.HandleFunc("GET /api/healthz", handlerReadiness)

	sMux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)

	sMux.HandleFunc("GET /api/reset", apiConfig.handlerReset)

	sMux.HandleFunc("POST /api/chirps", db.handlerPostChirps)

	sMux.HandleFunc("GET /api/chirps", db.handlerGetChirps)

	log.Printf("Serving files from %v on port: %v", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}