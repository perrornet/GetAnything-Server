package server

import (
	"fmt"
	"github.com/PerrorOne/GetAnything-Server/download"
	error2 "github.com/PerrorOne/GetAnything-Server/error"
	"github.com/PerrorOne/GetAnything-Server/extractors"
	"github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"time"
)

var Cache = cache.New(5*time.Minute, 10*time.Minute)

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
	if v, ok := Cache.Get(url); ok {
		resp := &Response{Code: 0, Data: v.(Data)}
		ctx.JSON(200, resp)
		log.Info(fmt.Sprintf("return Data from cahce key:%s", url))
		return
	}
	if url == "" {
		resp := &Response{Code: 1, Msg: error2.ClientError.Error()}
		ctx.JSON(400, resp)
		c, _ := ioutil.ReadAll(ctx.Request.Body)
		log.Info(string(c))
		log.Error(fmt.Sprintf("参数错误: %v|%s", ctx.Request.PostForm, ctx.Request.Body, string(c)))
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
	Cache.SetDefault(url, resp.Data)
	ctx.JSON(200, resp)
	log.Info(resp.String())
	return
}
