package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) LoginUser(email string, password string) (User, error) {
	user, err := db.GetUserByEmail(email)

	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		return User{}, errors.New("invalid password")
	}

	return *user, nil
}
