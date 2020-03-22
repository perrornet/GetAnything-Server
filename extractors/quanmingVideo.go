package extractors

import (
	"GetAnything-Server/download"
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
)

var (
	quanmingVideoInfo = regexp.MustCompile(`window\.metaInfo =(?s:(.*?));if`)
)

type quanmingResponse struct {
	VideoInfo struct {
		Vid        string `json:"vid"`
		ClarityURL []struct {
			URL string `json:"url"`
		} `json:"clarityUrl"`
	} `json:"videoInfo"`
}

type quanmingVideo struct {
	url    string
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (q *quanmingVideo) Init(url string) error {
	q.url = url
	q.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (q *quanmingVideo) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (q *quanmingVideo) GetFileInfo() ([]download.Info, error) {
	resp, err := q.client.Get(q.url, nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	data := quanmingVideoInfo.FindAllSubmatch(c, 1)
	if len(data) == 0 {
		return nil, errors.New("未能匹配到视频信息")
	}
	quanming := new(quanmingResponse)
	if err := json.Unmarshal(data[0][1], quanming); err != nil {
		return nil, err
	}
	if len(quanming.VideoInfo.ClarityURL) == 0 || quanming.VideoInfo.ClarityURL[0].URL == "" {
		return nil, errors.New("未能匹配到视频下载url")
	}

	return []download.Info{{Type: "mp4", Title: "全民小视频_" + quanming.VideoInfo.Vid, Url: quanming.VideoInfo.ClarityURL[0].URL}}, nil
}
