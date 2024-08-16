package main

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error){
	newDB := DB{path: "./database.json", mux: &sync.RWMutex{}}
	newDB.ensureDB()
	return &newDB, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	newChirp := Chirp{Body: body}
	currentDB, err := db.loadDB()
	if err != nil {
		return newChirp, err
	}

	chirps, err := db.GetChirps()
	if err != nil {
		return newChirp, err
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })

	if len(chirps) == 0 {
		newChirp.ID = 1
	} else {
		newChirp.ID = chirps[len(chirps) - 1].ID + 1
	}

	chirps = append(chirps, newChirp)

	for i, chirp := range chirps {
		currentDB.Chirps[i + 1] = chirp
	}

	db.writeDB(currentDB)
	return newChirp, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirpsSlice := []Chirp{}

	currentDB, err := db.loadDB()
	if err != nil {
		return chirpsSlice, err
	}

	dat, err := os.ReadFile(db.path)
	if err != nil {
		return chirpsSlice, err
	}
	
	err = json.Unmarshal(dat, &currentDB)
	if err != nil {
		return chirpsSlice, err
	}
	for _, chirp := range currentDB.Chirps {
		chirpsSlice = append(chirpsSlice, chirp)
	}
	return chirpsSlice, nil
}

// GetSingleChirp returns a chirp in the database using a chirpID
func (db *DB) GetSingleChirp(chirpID int) (Chirp, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	var chirp Chirp

	currentDB, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	dat, err := os.ReadFile(db.path)
	if err != nil {
		return chirp, err
	}
	
	err = json.Unmarshal(dat, &currentDB)
	if err != nil {
		return chirp, err
	}
	
	chirp, exist := currentDB.Chirps[chirpID]
	if !exist {
		return chirp, errors.New("Chirp you are looking for does not exist")
	}

	return chirp, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error{
	_, exist := os.Stat("./database.json")
	if exist != nil {
		newDBStructure := DBStructure{Chirps: map[int]Chirp{}}
		dat, err := json.Marshal(newDBStructure)
		if err != nil {
			return err
		}
		os.WriteFile("./database.json", dat, 0644)
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	currentDB := DBStructure{}
	dat, err := os.ReadFile("./database.json")
	if err != nil {
		return currentDB, err
	}
	err = json.Unmarshal(dat, &currentDB)
	if err != nil {
		return currentDB, err
	}
	return currentDB, nil
}


// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	os.WriteFile(db.path, dat, 0644)	
	return nil
}