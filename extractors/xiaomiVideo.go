package extractors

import (
	"GetAnything-Server/download"
	"encoding/json"
	"errors"
	"io/ioutil"
	url2 "net/url"
	"strings"
)

var (
	xiaomiVideoId    = "https://m.video.xiaomi.com/api/a3/play?id="
	xiaomiVideoUrl   = "https://openapiv2.yilan.tv/plat/play?id="
	xiaomiVideoTitle = "https://openapiv2.yilan.tv/plat/videodetail?id="
)

type xiaomiVideo struct {
	url    *url2.URL
	client *download.Http
}

type xiaomiGetVidResponse struct {
	Msg      string `json:"msg"`
	PlayInfo []struct {
		PlayURL string `json:"play_url"`
	} `json:"play_info"`
}

type xiaomiVideoResp struct {
	Retcode  string `json:"retcode"`
	Bitrates []struct {
		URI  string `json:"uri"`
		Code string `json:"code"`
	} `json:"bitrates"`
}

type xiaomiVideoTitleResponse struct {
	Retcode string `json:"retcode"`
	Name    string `json:"name"`
}

func (x *xiaomiVideo) GetVid() (string, error) {
	var shareId string
	switch x.url.Host {
	case "img.cdn.mvideo.xiaomi.com":

		d := strings.Split(x.url.String(), "/")
		shareId = d[len(d)-1]
	default:

		shareIds := x.url.Query().Get("seg")
		if shareIds == "" || len(strings.Split(shareIds, "/")) < 2 {
			return shareIds, errors.New("未能在url中寻找到分享ID")
		}
		d := strings.Split(shareIds, "/")
		shareId = d[len(d)-1]
	}
	resp, err := x.client.Get(xiaomiVideoId+shareId, nil)
	if err != nil {
		return "", err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	xr := new(xiaomiGetVidResponse)
	if err := json.Unmarshal(c, xr); err != nil {
		return "", err
	}
	if xr.Msg != "OK" || len(xr.PlayInfo) == 0 {
		return "", errors.New("小米获取视频ID接口返回失败")
	}
	return xr.PlayInfo[0].PlayURL, nil
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (x *xiaomiVideo) Init(url string) error {
	x.url, _ = url2.Parse(url)
	x.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (x *xiaomiVideo) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (x *xiaomiVideo) GetFileInfo() ([]download.Info, error) {
	vid, err := x.GetVid()
	if err != nil {
		return nil, err
	}
	title := make(chan string)
	//todo: 协程池
	go func(url string, title chan string) {
		resp, err := x.client.Get(url, nil)
		if err != nil {
			title <- ""
			return
		}
		c, _ := ioutil.ReadAll(resp.Body)
		xt := new(xiaomiVideoTitleResponse)
		if err := json.Unmarshal(c, xt); err != nil {
			title <- ""
			return
		}
		title <- xt.Name
	}(xiaomiVideoTitle+vid, title)
	resp, err := x.client.Get(xiaomiVideoUrl+vid, nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	xr := new(xiaomiVideoResp)
	if err := json.Unmarshal(c, xr); err != nil {
		return nil, err
	}
	if xr.Retcode != "200" || len(xr.Retcode) == 0 {
		return nil, errors.New("小米视频详情接口返回失败")
	}
	// 寻找清晰度最高的链接
	var ld, sd, hd, url string
	for _, v := range xr.Bitrates {
		switch v.Code {
		case "ld":
			ld = v.URI
		case "sd":
			sd = v.URI
		case "hd":
			hd = v.URI
		}
	}
	if hd != "" {
		url = hd
	} else if sd != "" {
		url = sd
	} else if ld != "" {
		url = ld
	} else {
		return nil, errors.New("小米视频详情接口返回失败")
	}

	t := <-title
	return []download.Info{{Type: "mp4", Url: url, Title: t}}, nil
}
