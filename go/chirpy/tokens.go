package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	authToken = string
	userId    = int
)

const (
	issuer          = "chirpy"
	oneDayInSeconds = 86_400
)

func (api *apiConfig) CreateAuthToken(userId int, expiresInSeconds int) (authToken, error) {
	if expiresInSeconds <= 0 || expiresInSeconds > oneDayInSeconds {
		expiresInSeconds = oneDayInSeconds
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiresInSeconds))),
		Subject:   fmt.Sprint(userId),
	})

	token, err := unsignedToken.SignedString([]byte(api.jwtSecret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (api *apiConfig) Validate(tokenString authToken) (userId, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(api.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	subject := token.Claims.(*jwt.RegisteredClaims).Subject
	userId, err := strconv.Atoi(subject)

	if err != nil {
		return 0, err
	}

	return userId, nil
}
