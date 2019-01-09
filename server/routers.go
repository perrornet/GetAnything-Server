package server

import (
	"GetAnything-Server/logger"
	logger2 "github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
)

func routers() *gin.Engine {
	router := gin.Default()
	//router.Static("/static", root)
	log := logger.NewLogger("GetAnything", logger2.InfoLevel)
	router.Use(func(ctx *gin.Context) {
		ctx.Set("log", log)
	})
	router.POST("/GetVideoUrl", GetVideoUrl)
	return router
}
