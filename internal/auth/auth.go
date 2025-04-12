package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func CreateCookie(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

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
