package web

import (
	"encoding/base64"
	"go-blogger/src/database"
	"go-blogger/src/utils"
	"html/template"
	"mime"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Go Web App",
	})
}

func handleFile(c *gin.Context) {
	fileName := c.Param("file_name")

	fileData := database.GetFileByName(fileName)

	if len(fileData) == 0 {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	// Set the appropriate response headers
	extension := regexp.MustCompile(`\.[a-zA-Z0-9]+$`).FindString(fileName)
	c.Data(http.StatusOK, mime.TypeByExtension(extension), fileData)
}

func needsAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		incomingSecret := c.GetHeader("Authorization")
		if incomingSecret != "Bearer "+secretKey {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func handleUpload(c *gin.Context) {
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
}

func handleCreateBlog(c *gin.Context) {
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
}

func handleBlog(c *gin.Context) {
	blogID := c.Param("blog_id")
	title := c.Param("title")
	realTitle := utils.Slugify(c.GetString("blog_title"))

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
}

func blogExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		blogID := c.Param("blog_id")
		title := database.GetBlogTitle(blogID)

		if title == "" {
			c.String(http.StatusNotFound, "Blog not found")
			c.Abort()
			return
		}
		c.Set("blog_title", title)
		c.Next()
	}
}

func handleBlogRedirect(c *gin.Context) {
	blogID := c.Param("blog_id")
	title := utils.Slugify(c.GetString("blog_title"))

	c.Redirect(http.StatusMovedPermanently, "/blog/"+blogID+"/"+title)
}
