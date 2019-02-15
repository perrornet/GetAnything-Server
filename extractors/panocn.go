package extractors

import (
	"bytes"
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
	"io/ioutil"
	url2 "net/url"
	"regexp"
)

var (
	panocnVideo = regexp.MustCompile(`"url_list":\["(.*?)",`)
)

type panocn struct {
	url    *url2.URL
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (p *panocn) Init(url string) error {
	p.url, _ = url2.Parse(url)
	p.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (p *panocn) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (p *panocn) GetFileInfo() ([]download.Info, error) {
	resp, err := p.client.Get(p.url.String(), nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	urls := panocnVideo.FindAll(c, 1)
	if len(urls) == 0 {
		return nil, errors.New("未能在火山小视频中解析到视频下载url")
	}
	urls = bytes.Split(urls[0], []byte(`"`))
	url := urls[len(urls)-2]
	url = bytes.Replace(url, []byte(`\u0026`), []byte("&"), -1)
	url = bytes.Replace(url, []byte(`\`), []byte(""), -1)
	u, err := url2.Parse(string(url))
	if err != nil {
		return nil, errors.New("未能在火山小视频中解析到视频下载url")
	}
	return []download.Info{{Title: u.Query().Get("video_id"), Url: u.String(), Type: "mp4"}}, nil
}
