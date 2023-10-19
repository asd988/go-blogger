package main

import (
	"go-blogger/src/database"
	"go-blogger/src/secret"
	"go-blogger/src/web"
)

func main() {
	database.InitDB()
	secretKey := secret.GetSecret()

	web.RunServer(secretKey)
}
