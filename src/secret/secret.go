package secret

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
)

const secretFilePath = "secret"

func generateSecret() string {
	randomBytes := make([]byte, 64)

	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err)
	}

	secret := base64.StdEncoding.EncodeToString(randomBytes)
	return secret
}

func GetSecret() string {
	bytes, err := os.ReadFile(secretFilePath)

	var secret string

	if err != nil {
		secret = generateSecret()

		err := os.WriteFile(secretFilePath, []byte(secret), 0644)
		if err != nil {
			log.Fatal(err)
		}

		println("Secret file created.")
	} else {
		secret = string(bytes)
	}

	return secret
}
