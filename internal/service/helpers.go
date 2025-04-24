package service

import "fmt"

func generateShortURL(baseURL, shortURLID string) string {
	return fmt.Sprintf("%s/%s", baseURL, shortURLID)
}
