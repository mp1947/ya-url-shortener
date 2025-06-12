package service

import "fmt"

// generateShortURL constructs a full short URL by combining the provided base URL and short URL identifier.
// It returns the resulting URL as a string in the format: baseURL/shortURLID.
func generateShortURL(baseURL, shortURLID string) string {
	return fmt.Sprintf("%s/%s", baseURL, shortURLID)
}
