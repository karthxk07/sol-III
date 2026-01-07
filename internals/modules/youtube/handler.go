package youtube

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func Register(r *gin.Engine) {
	youtube_router := r.Group("/youtube")
	youtube_router.POST("/upload", func(c *gin.Context) {

		title := c.PostForm("title")
		description := c.PostForm("description")
		categoryId, _ := strconv.Atoi(c.PostForm("categoryId"))
		tags := strings.Split(c.PostForm("tags"), ",")
		isReel := c.PostForm("isReel") == "true"

		snippet := Snippet{
			Title:       title,
			Description: description,
			Tags:        tags,
			CategoryID:  categoryId,
		}

		fileHeader, err := c.FormFile("video")
		if err != nil {
			c.JSON(400, gin.H{"error": "video file required"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		if isReel {
			newFile, newSize, err := ConvertToShort(file)
			if err != nil {
				c.JSON(500, gin.H{"error": "short conversion failed: " + err.Error()})
				return
			}
			file = newFile
			fileHeader.Size = newSize
		}

		if err := UploadYouTubeVideo(snippet, file, fileHeader.Size); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "uploaded"})
	})
}
