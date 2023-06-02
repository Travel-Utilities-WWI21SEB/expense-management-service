package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ContainsEmptyString(strings ...string) bool {
	for _, s := range strings {
		if s == "" {
			return true
		}
	}

	return false
}

func GenerateRandomString(length int) string {
	// Create a slice to store the generated characters
	result := make([]byte, length)

	// Calculate the length of the character set
	charsetLength := big.NewInt(int64(len(charset)))

	// Generate random index within the character set for each character in the result
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charsetLength)
		result[i] = charset[randomIndex.Int64()]
	}

	// Convert the byte slice to a string and return the result
	return string(result)
}
