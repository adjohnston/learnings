package db

import (
	"errors"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (db *DB) GetUserByEmail(email string) (*User, error) {
	dbStructure, err := db.load()
	if err != nil {
		return nil, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (db *DB) getUserById(id int) (*User, error) {
	dbStructure, err := db.load()
	if err != nil {
		return nil, err
	}

	for _, user := range dbStructure.Users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (db *DB) CreateUser(email string, password []byte) (User, error) {
	dbStructure, err := db.load()
	if err != nil {
		return User{}, err
	}

	existingUser, _ := db.GetUserByEmail(email)
	if existingUser != nil {
		return User{}, errors.New("user already exists")
	}

	id := len(dbStructure.Users) + 1
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	dbStructure.Users[id] = user
	err = db.write(dbStructure)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(id int, email string, hashedPassword []byte) (User, error) {
	dbStructure, err := db.load()
	if err != nil {
		return User{}, err
	}

	existingUser, err := db.getUserById(id)
	if err != nil {
		return User{}, errors.New("user already exists")
	}

	existingUser.Email = email
	existingUser.Password = hashedPassword
	dbStructure.Users[id] = *existingUser
	err = db.write(dbStructure)

	if err != nil {
		return User{}, err
	}

	return *existingUser, nil
}
