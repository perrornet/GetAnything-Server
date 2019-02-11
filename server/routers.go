package server

import (
	"github.com/PerrorOne/GetAnything-Server/logger"
	logger2 "github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func routers() *gin.Engine {
	router := gin.Default()
	//router.Static("/static", root)
	log := logger.NewLogger("GetAnything", logger2.InfoLevel)
	router.Use(func(ctx *gin.Context) {
		ctx.Set("log", log)
		method := ctx.Request.Method
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("charset", "utf-8")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		ctx.Next()
	})
	router.POST("/GetVideoUrl", GetVideoUrl)
	router.NoRoute(func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	return router
}
