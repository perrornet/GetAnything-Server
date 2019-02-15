package extractors

import (
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
	"io/ioutil"
	url2 "net/url"
	"regexp"
)

var (
	haokanVideo   = regexp.MustCompile(`video/mp4" src="(?s:(.*?))"`)
	haokanTitle   = regexp.MustCompile(`name="description" content="(?s:(.*?))"`)
	haokanHeaders = map[string]string{"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.81 Safari/537.36"}
)

type haokan struct {
	url    *url2.URL
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (h *haokan) Init(url string) error {
	h.url, _ = url2.Parse(url)
	h.client = download.NewHttp(haokanHeaders)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (h *haokan) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (h *haokan) GetFileInfo() ([]download.Info, error) {
	resp, err := h.client.Get(h.url.String(), nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	result := haokanVideo.FindAllSubmatch(c, 1)
	if len(result) == 0 {
		return nil, errors.New("好看视频url匹配失败")
	}
	var title string
	titleResult := haokanTitle.FindAllSubmatch(c, 1)
	if len(titleResult) != 0 { // 06d715b4ff41c2e2dd0b93272a1b7a05
		title = string(titleResult[0][1])
	}
	return []download.Info{{Title: title, Url: string(result[0][1]), Type: "mp4"}}, nil
}
