package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	md, err := os.ReadFile("test.md")
	if err != nil {
		panic(err)
	}

	p := parser.NewWithExtensions(parser.CommonExtensions | parser.HardLineBreak)
	html := markdown.ToHTML(md, p, nil)
	println(string(html))

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
			"BlogContent": template.HTML(string(html)),
		})
	})

	r.Run("127.0.0.1:8080") // Run the server on port 8080
}
