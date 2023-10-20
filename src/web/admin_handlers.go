package web

import (
	"encoding/base64"
	"go-blogger/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func handleGetFileHashes(c *gin.Context) {
	var json struct {
		BlogId string `json:"blog_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	snapshot_id := database.GetBlogSnapshotId(json.BlogId)
	files := database.GetSnapshotFiles(snapshot_id)

	var hashes []string
	for _, file := range files {
		hashes = append(hashes, base64.StdEncoding.EncodeToString(file.Hash))
	}

	c.JSON(http.StatusOK, gin.H{
		"file_hashes": hashes,
	})
}

func handleDownload(c *gin.Context) {
	var json struct {
		Hash string `json:"hash" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	hashBytes, err := base64.StdEncoding.DecodeString(json.Hash)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	name, data := database.GetFile(hashBytes)
	if data == nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+name)
	respondWithFile(c, name, data)
}
