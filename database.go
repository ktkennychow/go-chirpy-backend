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

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	HashedPassword []byte `json:"hashedpassword"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error){
	newDB := DB{path: path, mux: &sync.RWMutex{}}
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

	chirps, err := db.ReadChirps()
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

// ReadChirps returns all chirps in the database
func (db *DB) ReadChirps() ([]Chirp, error){
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

// ReadSingleChirp returns a chirp in the database using a chirpID
func (db *DB) ReadSingleChirp(chirpID int) (Chirp, error){
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

// CreateUsers creates a new User and saves it to disk
func (db *DB) CreateUsers(email string, hashedPassword []byte) (User, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	newUser := User{}

	currentDB, err := db.loadDB()
	if err != nil {
		return newUser, err
	}

	users, err := db.ReadUsers()
	if err != nil {
		return newUser, err
	}

	for _, v := range users {
		if v.Email == newUser.Email {
			return newUser, errors.New("a user with the input email already exists")
		}
	}

	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })

	if len(users) == 0 {
		newUser.ID = 1
	} else {
		newUser.ID = users[len(users) - 1].ID + 1
	}

	newUser.Email = email
	newUser.HashedPassword = hashedPassword

	users = append(users, newUser)

	for i, user := range users {
		currentDB.Users[i + 1] = user
	}

	db.writeDB(currentDB)
	return newUser, nil
}

// UpdateUser update a User and saves it to disk
func (db *DB) UpdateUser(email string, hashedPassword []byte, userID int) (User, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	updatedUser := User{}

	currentDB, err := db.loadDB()
	if err != nil {
		return updatedUser, err
	}
	
	user, exist := currentDB.Users[userID]
	if !exist {
		return updatedUser, errors.New("Chirp you are looking for does not exist")
	}

	updatedUser.ID = user.ID
	updatedUser.Email = email
	updatedUser.HashedPassword = hashedPassword

	currentDB.Users[userID] = updatedUser

	db.writeDB(currentDB)
	return updatedUser, nil
}

// ReadUsers returns all users in the database
func (db *DB) ReadUsers() ([]User, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	usersSlice := []User{}

	currentDB, err := db.loadDB()
	if err != nil {
		return usersSlice, err
	}

	for _, user := range currentDB.Users {
		usersSlice = append(usersSlice, user)
	}

	return usersSlice, nil
}

// ReadSingleUserbyEmail returns a user in the database
func (db *DB) ReadSingleUserbyEmail(userEmail string) (User, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	currentDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range currentDB.Users {
		if user.Email == userEmail {
			return user, nil
		}
	}

	return User{}, errors.New("no user with a matching email")
}

// ReadSingleUserbyID returns a user in the database
func (db *DB) ReadSingleUserbyID(userEmail string) (User, error){
	db.mux.RLock()
	defer db.mux.RUnlock()

	currentDB, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	dat, err := os.ReadFile(db.path)
	if err != nil {
		return User{}, err
	}
	
	err = json.Unmarshal(dat, &currentDB)
	if err != nil {
		return User{}, err
	}

	for _, user := range currentDB.Users {
		if user.Email == userEmail {
			return user, nil
		}
	}

	return User{}, errors.New("no user with a matching email")
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error{
	_, exist := os.Stat(db.path)
	if exist != nil {
		newDBStructure := DBStructure{Chirps: map[int]Chirp{}, Users: map[int]User{}}
		dat, err := json.Marshal(newDBStructure)
		if err != nil {
			return err
		}
		os.WriteFile(db.path, dat, 0644)
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	currentDB := DBStructure{}
	dat, err := os.ReadFile(db.path)
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
	err = os.WriteFile(db.path, dat, 0644)	
	if err != nil {
		return err
	}
	return nil
}

func deleteDB(path string) error{
	err := os.Remove(path)
	if err != nil {
		return err
	}
	
	return nil
}