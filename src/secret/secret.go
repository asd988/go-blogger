package secret

import (
	"encoding/base64"
	"go-blogger/src/genrandom"
	"log"
	"os"
)

const secretFilePath = "secret"

func generateSecret() string {
	randomBytes := genrandom.GenerateRandomBytes(64)
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
