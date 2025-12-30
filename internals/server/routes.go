package server

import (
	"github.com/gin-gonic/gin"
	"github.com/karthxk07/sol-III/internals/modules/google"
)

func Register(r *gin.Engine) {

	//register differnt routes ; add new routes to list, and respective Register func in the handler.go file
	google.Register(r)
}

