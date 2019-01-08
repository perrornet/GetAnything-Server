package extractors

import (
	"GetAnything-Server/download"
	error2 "GetAnything-Server/error"
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

type douyuTV struct {
	content string
}

var (
	douyuVideoInfoUrl = "http://vmobile.douyu.com/video/getInfo?vid="
	douyuTitle        = regexp.MustCompile("<h1>(.+?)</h1>")
	douyuHeaders      = map[string]string{"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"}
)

type douyuResp struct {
	Error int         `json:"error"`
	Data  interface{} `json:"data"`
}

func (d *douyuTV) GetFileFormUrl(url string) (string, error) {
	if !strings.Contains(url, "show/") {
		return "", error2.UrlError
	}
	url = strings.Replace(url, "vmobile.douyu.com", "v.douyu.com", 1)
	t := strings.Split(url, "/")
	if len(t) < 5 {
		return "", error2.UrlError
	}
	vId := strings.Split(t[4], "?")[0]
	client := download.NewHttp(douyuHeaders, true)
	resp, err := client.Get(url, nil)
	if err != nil {
		return "", err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	d.content = string(c)
	resp, err = client.Get(douyuVideoInfoUrl+vId, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data := &douyuResp{}
	c, _ = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(c, data); err != nil {
		return "", err
	}
	if data.Error != 0 {
		return "", errors.New(data.Data.(string))
	}
	if v, ok := data.Data.(map[string]interface{}); ok {
		if v1, ok := v["video_url"]; ok {
			return v1.(string), nil
		}
		return "", errors.New("请求斗鱼视频信息获取失败，data.data中没有'video_url'键")
	}
	return "", errors.New("斗鱼获取视频信息接口类型变化")
}

func (d *douyuTV) GetFileInfo() *download.Info {
	if d.content == "" {
		return nil
	}
	titles := douyuTitle.FindAllString(d.content, 1)
	for _, t := range titles {
		t = strings.Replace(t, "<h1>", "", 1)
		t = strings.Replace(t, "</h1>", "", 1)
		return &download.Info{Title: t, DownloadHeaders: douyuHeaders}
	}
	return &download.Info{DownloadHeaders: douyuHeaders}
}
