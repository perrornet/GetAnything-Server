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
	baiduHeaders    = map[string]string{"User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.81 Mobile Safari/537.36"}
	baiduVideoUrl   = regexp.MustCompile(`<video  src="(?s:(.*?))"`)
	baiduVideoTitle = regexp.MustCompile(`title: '(?s:(.*?))',`)
)

type baidu struct {
	url    *url2.URL
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (b *baidu) Init(url string) error {
	b.url, _ = url2.Parse(url)
	b.client = download.NewHttp(baiduHeaders)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (b *baidu) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (b *baidu) GetFileInfo() ([]download.Info, error) {
	resp, err := b.client.Get(b.url.String(), nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	content := baiduVideoUrl.FindAllSubmatch(c, 1)
	if len(content) == 0 {
		return nil, errors.New("未能匹配到视频下载url")
	}
	url := bytes.Replace(content[0][1], []byte("amp;"), []byte(""), -1)
	titles := baiduVideoTitle.FindAllSubmatch(c, 1)
	var title string
	if len(titles) != 0 {
		title = string(titles[0][1])
	}
	return []download.Info{{Url: string(url), Type: "mp4", Title: title}}, nil
}
