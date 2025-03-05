package usecase

import (
	"crypto/sha256"
	"fmt"
)

const resultStrLength = 8

func GenerateRandomIDFromURL(url string) string {
	hash := sha256.New()

	hash.Write([]byte(url))

	result := hash.Sum(nil)

	if len(result) < resultStrLength {
		return fmt.Sprintf("%x", result)
	}

	return fmt.Sprintf("%x", result)[:resultStrLength]
}
