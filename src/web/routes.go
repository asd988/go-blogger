package web

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	r.GET("/", handleIndex)
	r.GET("/blog/:blog_id", handleBlogRedirect)
	r.GET("/blog/:blog_id/:title", handleBlog)
	r.GET("/file/:file_name", handleFile)

	r.POST("/upload", handleUpload)
	r.POST("/create_blog", handleCreateBlog)
}
