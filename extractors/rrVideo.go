package extractors

import (
	"GetAnything-Server/download"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	url2 "net/url"
	"regexp"
)

var (
	rrVideoApi     = "https://api.rr.tv/v3plus/video/getVideoPlayLinkByVideoId?videoId=%s"
	rrVideoTitleRe = regexp.MustCompile(`tv/(?s:(.*?))\.mp4`)
	rrVideoId      = regexp.MustCompile(`id=(?s:(\d+))`)
	rrVideoHeaders = map[string]string{
		"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.81 Safari/537.36",
		"clientVersion": "undefined",
		"clientType":    "web",
		"token":         "undefined",
	}
)

type rrVideoResponse struct {
	Code string `json:"code"`
	Data struct {
		PlayLink string `json:"playLink"`
	} `json:"data"`
}

type rrVideo struct {
	url    *url2.URL
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (r *rrVideo) Init(url string) error {
	r.url, _ = url2.Parse(url)
	r.client = download.NewHttp(rrVideoHeaders)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (r *rrVideo) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (r *rrVideo) GetFileInfo() ([]download.Info, error) {
	data := rrVideoId.FindAllSubmatch([]byte(r.url.String()), 1)
	if len(data) == 0 {
		return nil, errors.New("url错误,未能寻找到id")
	}
	resp, err := r.client.Get(fmt.Sprintf(rrVideoApi, string(data[0][1])), nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	rr := new(rrVideoResponse)
	if err := json.Unmarshal(c, rr); err != nil {
		return nil, err
	}
	if rr.Code != "0000" {
		return nil, errors.New("当前视频类型不支持")
	}
	if rr.Data.PlayLink == "" {
		return nil, errors.New("当前视频类型不支持")
	}
	data = rrVideoTitleRe.FindAllSubmatch([]byte(rr.Data.PlayLink), 1)
	if len(data) != 0 {
		return []download.Info{{Title: string(data[0][1]), Url: rr.Data.PlayLink, Type: "mp4"}}, nil
	}
	return []download.Info{{Url: rr.Data.PlayLink, Type: "mp4"}}, nil
}
