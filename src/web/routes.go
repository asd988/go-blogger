package web

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	r.GET("/", handleIndex)
	r.GET("/file/:file_name", handleFile)

	r.GET("/blog_list", handleBlogList)
	r.GET("/blog/:blog_id", blogExists(), handleBlogRedirect)
	r.GET("/blog/:blog_id/:first", blogExists(), handleBlog)
	r.GET("/blog/:blog_id/:first/:second", blogExists(), handleSnapshotBlog)

	r.POST("/upload", needsAuth(), handleUpload)
	r.POST("/create_blog", needsAuth(), handleCreateBlog)
	r.POST("/get_file_hashes", needsAuth(), handleGetFileHashes)
	r.POST("/download", needsAuth(), handleDownload)
}
