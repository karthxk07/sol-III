package google

import (
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {

	//bind route /google to the main engine
	google_router := r.Group("/google")

	//GET the oauth url
	google_router.GET("/getOAuthurl", func(c *gin.Context) {
		c.JSON(200, gin.H{"oauth_url": GetOAuthUrl()})
	})

	//handle the google oauth redirect
	google_router.GET("/code_redirect", func(c *gin.Context) {
		//parse the auth_code
		code, err := GetAuthCode(c.Request.URL)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "access_denied"})
		}

		//return a acess granted message
		c.JSON(200, gin.H{"message": "access granted", "code": code})

		//get the access token and the refresh token (shift this to the pipeline later)
		GetAccRefTokens(code)
	})
}
