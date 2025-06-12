package usecase

import (
	"crypto/sha256"
	"fmt"
)

const resultStrLength = 8

// GenerateIDFromURL generates a unique string identifier from the given URL using SHA-256 hashing.
// The resulting ID is a hexadecimal string truncated to resultStrLength characters if necessary.
// If the hash output is shorter than resultStrLength, the full hash is returned as a hex string.
//
// Parameters:
//   - url: The input URL to generate an ID for.
//
// Returns:
//   - A string representing the unique ID derived from the URL.
func GenerateIDFromURL(url string) string {
	hash := sha256.New()

	hash.Write([]byte(url))

	result := hash.Sum(nil)

	if len(result) < resultStrLength {
		return fmt.Sprintf("%x", result)
	}

	return fmt.Sprintf("%x", result)[:resultStrLength]
}
