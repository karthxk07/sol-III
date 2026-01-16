package gemini

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	gemini_router := r.Group("gemini")

	gemini_router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
}
