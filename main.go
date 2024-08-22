package main

import (
	"flag"
	"log"
	"net/http"
)

type apiConfig struct {
	FileserverHits int
}

func main(){
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
			deleteDB("./database.json")
	}

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

	sMux.HandleFunc("POST /api/chirps", db.handlerCreateChirps)

	sMux.HandleFunc("GET /api/chirps", db.handlerReadChirps)

	sMux.HandleFunc("GET /api/chirps/{chirpID}", db.handlerReadSingleChirp)

	sMux.HandleFunc("POST /api/users", db.handlerCreateUsers)

	sMux.HandleFunc("POST /api/login", db.handlerLogin)

	log.Printf("Serving files from %v on port: %v", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}