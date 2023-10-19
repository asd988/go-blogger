package main

import (
	"fmt"
	"mime"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	"go-blogger/src/database"
	secret "go-blogger/src/secret"
)

var secret_key string

func main() {
	database.InitDB()
	secret_key = secret.GetSecret()

	r := gin.Default()

	// Serve the static files (e.g., CSS, JavaScript, images) from the "static" directory
	r.Static("/static", "./static")

	// Define a route to render the HTML template
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Go Web App",
		})
	})
	r.GET("/blog", func(c *gin.Context) {
		c.HTML(http.StatusOK, "blog.html", gin.H{
			"title":       "Go Web App",
			"BlogContent": "template.HTML(string(html))",
		})
	})

	r.GET("/file/:file_name", func(c *gin.Context) {
		// Get the file name parameter from the URL
		fileName := c.Param("file_name")

		// Use the GetFile function to retrieve the file data
		fileData := database.GetFileByName(fileName)

		if len(fileData) == 0 {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		// Set the appropriate response headers
		// c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Data(http.StatusOK, mime.TypeByExtension(regexp.MustCompile(`\.[a-zA-Z0-9]+$`).FindString(fileName)), fileData)
	})

	r.POST("/upload", func(c *gin.Context) {
		incoming_secret := c.GetHeader("Authorization")
		if incoming_secret != "Bearer "+secret_key {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		defer file.Close()

		name := header.Filename
		println("Name ", name)

		// Read the file data into a byte slice
		data := make([]byte, header.Size)
		_, err = file.Read(data)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading the file")
			return
		}

		p := parser.NewWithExtensions(parser.CommonExtensions | parser.HardLineBreak)
		html := markdown.ToHTML(data, p, nil)
		println(string(html))

		hash := database.StoreFile(name, data)
		fmt.Printf("Hash: %x\n", hash)

		c.String(http.StatusOK, "File uploaded and stored successfully")
	})

	r.Run("127.0.0.1:8080") // Run the server
}
