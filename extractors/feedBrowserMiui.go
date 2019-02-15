package extractors

import (
	"encoding/json"
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
	"io/ioutil"
	url2 "net/url"
	"strings"
)

var feedBrowserMiuiVideoInfoUrl = `https://feed.browser.miui.com/brec-api/api/v1/news/detail`

type feedBrowserMiuiResponse struct {
	ReturnCode int `json:"returnCode"`
	Data       []struct {
		Document struct {
			VideoURL string `json:"videoUrl"`
			Title    string `json:"title"`
		} `json:"document"`
	} `json:"data"`
}

type feedBrowserMiui struct {
	url    *url2.URL
	client *download.Http
	body   map[string]string
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (f *feedBrowserMiui) Init(url string) error {
	f.url, _ = url2.Parse(strings.Replace(url, "#", "?", 1))
	f.client = download.NewHttp(nil)
	f.body = map[string]string{
		"type":       "video",
		"docId":      "",
		"parameters": "%7B%22appName%22%3A%22mifeeds%22%2C%22traceId%22%3A%22%22%2C%22ref%22%3A%22%22%7D",
	}
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (f *feedBrowserMiui) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (f *feedBrowserMiui) GetFileInfo() ([]download.Info, error) {
	f.body["docId"] = f.url.Query().Get("docid")
	if f.body["docId"] == "" {
		return nil, errors.New("url错误,未能在url中寻找到docId")
	}
	resp, err := f.client.Post(feedBrowserMiuiVideoInfoUrl, nil, f.body)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	fr := &feedBrowserMiuiResponse{}
	if err := json.Unmarshal(c, fr); err != nil {
		return nil, err
	}
	if fr.ReturnCode != 200 || len(fr.Data) == 0 {
		return nil, errors.New("请求接口出现错误")
	}
	return []download.Info{{Type: "mp4", Title: fr.Data[0].Document.Title, Url: fr.Data[0].Document.VideoURL}}, nil
}
