package genrandom

import (
	"crypto/rand"
	"log"
)

func GenerateRandomBytes(length int) []byte {
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err)
	}

	return randomBytes
}

func GenerateRandomString(length int) string {
	randomBytes := GenerateRandomBytes(length)

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	result := ""
	for _, b := range randomBytes {
		result += string(letters[int(b)%0b111111])
	}

	return result
}
