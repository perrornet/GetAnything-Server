package server

import (
	"github.com/gin-gonic/gin"
)

func routers(root string)*gin.Engine{
	router := gin.Default()
	router.Static("/static", root)
	router.GET("/GetVideoUrl", func(c *gin.Context) {
		
	})
	return router
}
