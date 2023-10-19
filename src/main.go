package main

import (
	"encoding/base64"
	"html/template"
	"mime"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	"go-blogger/src/database"
	"go-blogger/src/secret"
)

var secretKey string

func slugify(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Remove non-alphanumeric and non-hyphen characters
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove any leading or trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

func main() {
	database.InitDB()
	secretKey = secret.GetSecret()

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
		incomingSecret := c.GetHeader("Authorization")
		if incomingSecret != "Bearer "+secretKey {
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

		hash := database.StoreFile(name, data)
		resp := base64.StdEncoding.EncodeToString(hash)

		c.JSON(http.StatusOK, gin.H{
			"hash": resp,
		})
	})

	r.POST("/create_blog", func(c *gin.Context) {
		incomingSecret := c.GetHeader("Authorization")
		if incomingSecret != "Bearer "+secretKey {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		var json struct {
			Title      string   `json:"title" binding:"required"`
			PageHash   string   `json:"page_hash" binding:"required"`
			FileHashes []string `json:"file_hashes" binding:"required"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}

		println("Creating blog with title ", json.Title)

		pageHash, err := base64.StdEncoding.DecodeString(json.PageHash)
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}

		var hashes [][]byte
		for _, hash := range json.FileHashes {
			hashBytes, err := base64.StdEncoding.DecodeString(hash)
			if err != nil {
				c.String(http.StatusBadRequest, "Bad request")
				return
			}
			hashes = append(hashes, hashBytes)
		}

		id := database.CreateBlog(json.Title, pageHash, hashes)

		c.JSON(http.StatusOK, gin.H{
			"blog_id": id,
		})
	})

	r.GET("/blog/:blog_id", func(c *gin.Context) {
		blogID := c.Param("blog_id")
		title := slugify(database.GetBlogTitle(blogID))

		c.Redirect(http.StatusMovedPermanently, "/blog/"+blogID+"/"+title)
	})

	r.GET("/blog/:blog_id/:title", func(c *gin.Context) {
		blogID := c.Param("blog_id")
		title := c.Param("title")
		realTitle := slugify(database.GetBlogTitle(blogID))

		if title != realTitle {
			c.String(http.StatusNotFound, "Blog not found")
			return
		}

		data := database.GetBlogContent(blogID)
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.HardLineBreak)
		html := markdown.ToHTML(data, p, nil)

		c.HTML(http.StatusOK, "blog.html", gin.H{
			"title":       "Go Web App",
			"BlogTitle":   title,
			"BlogContent": template.HTML(string(html)),
		})
	})

	r.Run("127.0.0.1:8080") // Run the server
}
