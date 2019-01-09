package server

import (
	"GetAnything-Server/download"
	error2 "GetAnything-Server/error"
	"GetAnything-Server/extractors"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
)

type Data struct {
	Headers map[string]string `json:"headers"`
	Info    []download.Info   `json:"info"`
}
type Response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data Data   `json:"data"`
}

func (r *Response) String() string {
	return fmt.Sprintf("<code=%d, data.headers=%v, data.Info=%v>", r.Code, r.Data.Headers, r.Data.Info)
}

func GetVideoUrl(ctx *gin.Context) {
	l, _ := ctx.Get("log")
	log := l.(*logger.Logger)
	url := ctx.PostForm("url")
	if url == "" {
		resp := &Response{Code: 1, Msg: error2.ClientError.Error()}
		ctx.JSON(400, resp)
		log.Error("参数错误")
		return
	}
	log.Info(url)
	d, err := extractors.Match(url)
	if err != nil {
		resp := &Response{Code: 1, Msg: error2.NotFoundHandlerError.Error()}
		ctx.JSON(500, resp)
		log.Error(err.Error())
		return
	}

	info, err := d.GetFileInfo()
	if err != nil {
		log.Error(err.Error())
		resp := &Response{Code: 1, Msg: error2.ServerError.Error()}
		ctx.JSON(500, resp)
		return
	}
	resp := &Response{Code: 0, Data: Data{Headers: d.GetDownloadHeaders(), Info: info}}
	ctx.JSON(200, resp)
	log.Info(resp.String())
	return
}
