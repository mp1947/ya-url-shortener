// Package auth provides JWT-based authentication utilities for generating and validating tokens with user UUIDs.
package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Claims holds JWT claims, embedding standard fields and a user ID.
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))

// CreateToken generates a JWT for the given userID using HS256 signing.
// Returns the token string or an error.
func CreateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate checks the validity of a JWT token string and returns whether it's valid and the associated user UUID.
// Returns false and uuid.Nil if the token is invalid or empty.
func Validate(tokenString string) (bool, uuid.UUID) {

	if tokenString == "" {
		return false, uuid.Nil
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secretKey, nil
		})
	if err != nil {
		return false, uuid.Nil
	}

	if !token.Valid {
		fmt.Println("token not valid")
		return false, uuid.Nil
	}

	return true, claims.UserID
}
