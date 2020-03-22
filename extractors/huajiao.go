package extractors

import (
	"GetAnything-Server/download"
	"GetAnything-Server/utils"
	"bytes"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	huajiaoVideoUrl = regexp.MustCompile(`"mp4":"(?s:(.*?))"`)
)

// 这个写的好无聊啊
type huajiao struct {
	url    string
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (h *huajiao) Init(url string) error {
	h.url = url
	h.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (h *huajiao) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (h *huajiao) GetFileInfo() ([]download.Info, error) {
	resp, err := h.client.Get(h.url, nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	urls := huajiaoVideoUrl.FindAllSubmatch(c, 1)
	if bytes.Equal(urls[0][1], []byte("")) {
		return nil, errors.New("未能在页面中匹配到视频下载url")
	}
	m := utils.NewMd5(urls[0][1])
	return []download.Info{{Type: "mp4", Title: m.Encrypt(), Url: strings.Replace(string(urls[0][1]), `\`, "", -1)}}, nil
}
