package server

import (
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	//create a new router
	r := gin.New()

	//register all routes
	Register(r)

	//return the engine
	return r
}
