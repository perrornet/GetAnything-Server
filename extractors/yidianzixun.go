package extractors

import (
	"GetAnything-Server/download"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	yidianzixunVideoUrl   = regexp.MustCompile(`video src="(?s:(.*?\.mp4))"`)
	yidianzixunVideoTitle = regexp.MustCompile(`name="description" content="(?s:(.*?))"`)
)

type yidianzixun struct {
	url    string
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (y *yidianzixun) Init(url string) error {
	y.url = url
	y.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (y *yidianzixun) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (y *yidianzixun) GetFileInfo() ([]download.Info, error) {
	resp, err := y.client.Get(y.url, nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	urls := yidianzixunVideoUrl.FindAllSubmatch(c, 1)
	if len(urls) == 0 {
		return nil, errors.New("未能匹配到视频URL")
	}
	if string(urls[0][1]) == "" {
		return nil, errors.New("未能匹配到视频URL")
	}
	var title string
	titles := yidianzixunVideoTitle.FindAllSubmatch(c, 1)
	if len(titles) != 0 {
		title = string(titles[0][1])
	}
	if title == "" {
		t := strings.Split(string(urls[0][1]), "/")
		title = strings.Split(t[len(t)-1], ".")[0]
	}
	return []download.Info{{Title: title, Type: "mp4", Url: string(urls[0][1])}}, nil
}
