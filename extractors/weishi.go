package extractors

import (
	"GetAnything-Server/download"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	url2 "net/url"
	"strings"
)

type weishiResponse struct {
	Ret  int `json:"ret"`
	Data struct {
		Feeds []struct {
			FeedDesc string `json:"feed_desc"`
			VideoURL string `json:"video_url"`
		} `json:"feeds"`
	} `json:"data"`
}

type weishi struct {
	url    *url2.URL
	apiUrl *url2.URL
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (w *weishi) Init(url string) error {
	w.url, _ = url2.Parse(url)
	w.client = download.NewHttp(nil)
	w.apiUrl, _ = url2.Parse(`https://h5.weishi.qq.com/webapp/json/weishi/WSH5GetPlayPage`)
	return nil
}

func (w *weishi) buldApiUrl(feedid string) string {
	v := w.apiUrl.Query()
	v.Set("t", "0.35894579952784755")
	v.Set("g_tk", "")
	v.Set("feedid", feedid)
	v.Set("recommendtype", "0")
	v.Set("datalvl", "")
	v.Set("qua", "")
	v.Set("uin", "")
	v.Set("format", "json")
	v.Set("inCharset", "utf-8")
	v.Set("outCharset", "utf-8")
	w.apiUrl.RawQuery = v.Encode()
	return w.apiUrl.String()
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (w *weishi) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (w *weishi) GetFileInfo() ([]download.Info, error) {
	paths := strings.Split(w.url.EscapedPath(), "/")
	if len(paths) == 0 || paths[2] != "feed" {
		return nil, errors.New("url错误,未能在该url中寻找到视频ID")
	}
	resp, err := w.client.Get(w.buldApiUrl(paths[3]), nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	data := &weishiResponse{}
	if err := json.Unmarshal(c, data); err != nil {
		return nil, err
	}
	if data.Ret != 0 {
		return nil, errors.New(fmt.Sprintf("微视接口请求失败,返回码: %d", data.Ret))
	}
	if len(data.Data.Feeds) == 0 {
		return nil, errors.New(fmt.Sprintf("微视接口请求失败,未能返回视频数据"))
	}
	return []download.Info{{Title: data.Data.Feeds[0].FeedDesc, Url: data.Data.Feeds[0].VideoURL, Type: "mp4"}}, nil
}
