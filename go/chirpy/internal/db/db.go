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
