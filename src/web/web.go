package web

import "github.com/gin-gonic/gin"

var secretKey string

func RunServer(secret string) {
	secretKey = secret

	r := gin.Default()

	// Serve the static files (e.g., CSS, JavaScript, images) from the "static" directory
	r.Static("/static", "./static")

	// Define a route to render the HTML template
	r.LoadHTMLGlob("templates/*")

	setupRoutes(r)

	r.Run("127.0.0.1:8080") // Run the server
}
