package app

import "math/rand/v2"

// generateUrlId generates random string with fixed length from charSet.
func generateUrlId(length int) string {
	var result []byte

	for i := 0; i < length; i++ {
		result = append(result, charSet[rand.IntN(len(charSet)-1)])
	}

	return string(result)
}
