package utils

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecretKey = []byte("your_jwt_secret_key")

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		return "", errors.New("missing UID in token claims")
	}

	return uid, nil
}
