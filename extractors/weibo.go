package extractors

import (
	"GetAnything-Server/download"
	"errors"
	"io/ioutil"
	"log"
	url2 "net/url"
	"regexp"
)

type weibo struct {
	content string
	url     string
	client  *download.Http
}

var (
	videoInfoUrl = "https://m.weibo.cn/s/video/object?object_id=%s&mid=%s"
	videoUrl1    = regexp.MustCompile(`sources=\\"(.*?)\\"`)
	videoUrl2    = regexp.MustCompile(`sources="(.*?)"`)
)

func (w *weibo) Init(url string) error {
	w.url = url
	w.client = download.NewHttp(nil)
	return nil
}

func (w *weibo) GetDownloadHeaders() map[string]string { return nil }

func (w *weibo) GetFileInfo() ([]download.Info, error) {
	data := make([]download.Info, 0)
	u, _ := url2.Parse(w.url)
	log.Println(u.Host)
	if u.Host == "weibo.com" {
		client := download.NewHttp(nil)
		resp, err := client.Get(w.url, nil)
		if err != nil {
			return data, err
		}
		c, _ := ioutil.ReadAll(resp.Body)
		w.content = string(c)
		s := videoUrl2.FindAllString(w.content, 1)
		if len(s) == 0 {
			return data, errors.New("微博视频接口发生变化。")
		}
		//todo
	}
	return data, errors.New("当前类型暂未支持！")
}
