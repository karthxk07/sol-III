package gemini

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuoteRequest struct {
	Quote string `json:"quote" binding:"required"`
}

func Register(r *gin.Engine) {
	gemini_router := r.Group("gemini")

	gemini_router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	gemini_router.POST("/get_metadata", func(c *gin.Context) {
		var req QuoteRequest

		// Bind and validate the JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Quote is required",
			})
			return
		}

		// Process the quote using the Gemini function
		videoMetadata, err := ProcessQuoteWithGemini(req.Quote)
		if err != nil {
			log.Fatal(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error processing the quote",
			})
			return
		}

		// Send the VideoMetadata as response
		c.JSON(http.StatusOK, videoMetadata)
	})
}
