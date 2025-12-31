package google

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Register(r *gin.Engine) {

	//bind route /google to the main engine
	google_router := r.Group("/google")

	//GET the oauth url
	google_router.GET("/getOAuthurl", func(c *gin.Context) {
		//get the oauth url
		oauth_url, err := GetOAuthUrl()
		if err != nil {
			log.Fatal("error getting the oauth url", err.Error())
			c.AbortWithStatusJSON(400, gin.H{"error": "error getting the oauth url"})
		}

		//return a status ok
		c.JSON(200, gin.H{"oauth_url": oauth_url})
	})

	//handle the google oauth redirect
	google_router.GET("/code_redirect", func(c *gin.Context) {
		//parse the auth_code
		code, err := GetAuthCode(c.Request.URL)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "access_denied"})
		}

		//pipeline to execute after
		{
			//get the token
			var token Token
			err := LoadAccRefToken(&token)
			if err != nil {
				log.Fatal("error while getting the token", err.Error())
			}

			//get the access token and the refresh token
			err = GetAccRefTokens(code, &token)
			if err != nil {
				log.Fatal("error while getting the access and refresh token:", err.Error())
			}

			//goto entry
			err = Prepare(&token)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": "error while preparing the oauth pipeline"})
			}
		}

		//return a acess granted message
		c.JSON(200, gin.H{"message": "access granted", "code": code})
	})
}
