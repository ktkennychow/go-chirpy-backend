package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	FileserverHits int
	DB *DB
	jwtSecret string
}

func main(){
	godotenv.Load()
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
	
	apiConfig := apiConfig{FileserverHits: 0, DB: db, jwtSecret: os.Getenv("JWT_SECRET")}

	sMux := http.NewServeMux()

	server := &http.Server{Handler: sMux, Addr: ":" + port }
	
	dir := http.Dir(".")
	handlerfs := apiConfig.middlewareMetricsInc(http.FileServer(dir))

	sMux.Handle("GET /app/*", http.StripPrefix("/app", handlerfs))

	sMux.HandleFunc("GET /api/healthz", handlerReadiness)

	sMux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)

	sMux.HandleFunc("GET /api/reset", apiConfig.handlerReset)

	sMux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirps)

	sMux.HandleFunc("GET /api/chirps", apiConfig.handlerReadChirps)

	sMux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerReadSingleChirp)

	sMux.HandleFunc("POST /api/users", apiConfig.handlerCreateUsers)

	sMux.HandleFunc("PUT /api/users", apiConfig.handlerModifyUsers)

	sMux.HandleFunc("POST /api/login", apiConfig.handlerLogin)

	sMux.HandleFunc("POST /api/refresh", apiConfig.handlerRefreshAuth)

	sMux.HandleFunc("POST /api/revoke", apiConfig.handlerRevokeAuth)

	log.Printf("Serving files from %v on port: %v", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}