package usecase

import "math/rand/v2"

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomID(length int) string {
	var result []byte

	for i := 0; i < length; i++ {
		result = append(result, charSet[rand.IntN(len(charSet)-1)])
	}

	return string(result)
}
