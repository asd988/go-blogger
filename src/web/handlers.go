package web

import (
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

func respondWithFile(c *gin.Context, fileName string, data []byte) {
	// Set the appropriate response headers
	extension := regexp.MustCompile(`\.[a-zA-Z0-9]+$`).FindString(fileName)
	c.Data(http.StatusOK, mime.TypeByExtension(extension), data)
}

func handleFile(c *gin.Context) {
	fileName := c.Param("file_name")

	fileData := database.GetFileByName(fileName)

	if len(fileData) == 0 {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	respondWithFile(c, fileName, fileData)
}

func handleBlog(c *gin.Context) {
	blogID := c.Param("blog_id")
	title := c.Param("title")
	realTitle := c.GetString("blog_title")

	if title == utils.Slugify(realTitle) {
		data := database.GetBlogContent(blogID)
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.HardLineBreak)
		html := markdown.ToHTML(data, p, nil)

		c.HTML(http.StatusOK, "blog.html", gin.H{
			"title":       "Go Web App",
			"BlogTitle":   title,
			"BlogContent": template.HTML(string(html)),
		})
		return
	} else {
		snapshot_id := database.GetBlogSnapshotId(blogID)
		files := database.GetSnapshotFiles(snapshot_id)
		// if title is in file_names
		for _, file := range files {
			if title == file.Name {
				_, data := database.GetFile(file.Hash)
				respondWithFile(c, file.Name, data)
				return
			}
		}
	}

	c.String(http.StatusNotFound, "Blog not found")
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

func handleBlogList(c *gin.Context) {
	blogs := database.GetBlogs(0, 10)
	for _, blog := range blogs {
		println("title: ", blog.Title)
	}

	c.HTML(http.StatusOK, "blog_list.html", gin.H{
		"title": "Go Web App",
		"Blogs": blogs,
	})
}
