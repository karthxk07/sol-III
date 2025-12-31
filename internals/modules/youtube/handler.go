package youtube

import (
	"github.com/gin-gonic/gin"
)

func Register(g *gin.Engine) {
	youtube := g.Group("/youtube")

	youtube.POST("/upload", func(c *gin.Context) {

		if c.Request.ContentLength <= 0 {
			c.AbortWithStatusJSON(400, gin.H{"error": "some error occured"})
			return
		}

		err := UploadVideo(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "placeholder"})
	})
}

