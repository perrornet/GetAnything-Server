package extractors

import (
	"GetAnything-Server/download"
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
)

var (
	huyaApiUrl = "http://v-api-play.huya.com/?format=mp4%2Cm3u8&r=vhuyaplay%2Fvideo&vid="
	huyaVid    = regexp.MustCompile(`/(?s:(.*?))\.html`)
	huyaTitle  = regexp.MustCompile(`data-text='(?s:(.*?))'`)
)

type huyaResponse struct {
	Code   int `json:"code"`
	Result struct {
		Items []struct {
			Format    string `json:"format"`
			Transcode struct {
				Urls []string `json:"urls"`
			} `json:"transcode"`
		} `json:"items"`
	} `json:"result"`
}

type huya struct {
	url    string
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (h *huya) Init(url string) error {
	h.url = url
	h.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (h *huya) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (h *huya) GetFileInfo() ([]download.Info, error) {
	vids := huyaVid.FindAllSubmatch([]byte(h.url), 1)
	if len(vids) == 0 {
		return nil, errors.New("未能在url中寻找到VID")
	}

	resp, err := h.client.Get(huyaApiUrl+string(vids[0][1]), nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	data := &huyaResponse{}
	if err := json.Unmarshal(c, data); err != nil {
		return nil, err
	}
	if data.Code != 1 && len(data.Result.Items) == 0 {
		return nil, errors.New("虎牙接口返回出错!")
	}
	var title string
	titles := huyaTitle.FindAllSubmatch(c, 1)
	if len(titles) != 0 {
		title = string(titles[0][1])
	}
	for _, v := range data.Result.Items {
		if v.Format == "mp4" {
			if len(v.Transcode.Urls) == 0 {
				break
			}
			return []download.Info{{Url: v.Transcode.Urls[0], Type: "mp4", Title: title}}, nil
		}
	}
	return nil, nil
}
