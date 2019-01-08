package server

import (
	"github.com/gin-gonic/gin"
)

func routers() *gin.Engine {
	router := gin.Default()
	//router.Static("/static", root)
	router.POST("/GetVideoUrl", GetVideoUrl)
	return router
}
