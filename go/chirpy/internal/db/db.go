package db

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Users  map[int]User
	Chirps map[int]Chirp
}

func (db *DB) ensure() error {
	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		return db.create()
	}

	return err
}

func (db *DB) create() error {
	data := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}

	return db.write(data)
}

func (db *DB) load() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) write(data DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, d, 0600)
	if err != nil {
		return err
	}

	return nil
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensure()

	return db, err
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.load()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStructure, err := db.load()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}

	return chirp, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.load()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.write(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.load()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:    id,
		Email: email,
	}
	dbStructure.Users[id] = user
	err = db.write(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
