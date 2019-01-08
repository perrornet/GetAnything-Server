package server

import (
	"GetAnything-Server/download"
	error2 "GetAnything-Server/error"
	"GetAnything-Server/extractors"
	"github.com/gin-gonic/gin"
	"log"
)

type Data struct {
	Url  string        `json:"url"`
	Info download.Info `json:"info"`
}
type Response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data Data   `json:"data"`
}

func GetVideoUrl(ctx *gin.Context) {
	url := ctx.PostForm("url")
	if url == "" {
		resp := &Response{Code: 1, Msg: error2.ClientError.Error()}
		ctx.JSON(400, resp)
		log.Println("传入url参数为空！")
		return
	}
	d, err := extractors.Match(url)
	if err != nil {
		resp := &Response{Code: 1, Msg: err.Error()}
		ctx.JSON(500, resp)
		log.Println(err)
		return
	}
	downloadUrl, err := d.GetFileFormUrl(url)
	if err != nil {
		log.Println(err)
		resp := &Response{Code: 1, Msg: error2.ServerError.Error()}
		ctx.JSON(500, resp)
		return
	}
	info := d.GetFileInfo()
	resp := &Response{Code: 0, Data: Data{Url: downloadUrl, Info: *info}}
	ctx.JSON(200, resp)
	return
}
