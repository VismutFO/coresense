package jwtutils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func GetTokenFromHeader(tokenHeader string) (string, error) {
	if tokenHeader == "" {
		return "", errors.New("token header is empty")
	}

	// Check Bearer token format
	if !strings.HasPrefix(tokenHeader, "Bearer ") {
		return "", errors.New("token header is not bearer token")
	}

	return strings.TrimPrefix(tokenHeader, "Bearer "), nil
}

func GetClaims(tokenString string, jwtSecret []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
