package extractors

import (
	"GetAnything-Server/download"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	url2 "net/url"
	"regexp"
	"strings"
)

var (
	sohuMVid        = regexp.MustCompile(`/(?s:(\d+))\.shtml`)
	sohuTVid        = regexp.MustCompile(`v/(?s:(.*?))\.html`)
	sohuSuVid       = regexp.MustCompile(`^/\d+/\d+/.*?\.mp4$`)
	sohuVideoApi    = "https://my.tv.sohu.com/play/videonew.do?vid=%s&ver=&ssl=1"
	sohuNewVideoApi = "https://data.vod.itc.cn/ip?prod=rtb&new="
)

type sohuResponse struct {
	Data struct {
		TvName   string   `json:"tvName"`
		ClipsURL []string `json:"clipsURL"`
		Su       []string `json:"su"`
	} `json:"data"`
}

type sohuNewResponse struct {
	Servers []struct {
		URL string `json:"url"`
	} `json:"servers"`
}

type sohu struct {
	url    *url2.URL
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (s *sohu) Init(url string) error {
	s.url, _ = url2.Parse(url)
	s.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (s *sohu) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (s *sohu) GetFileInfo() ([]download.Info, error) {
	var vid string
	switch s.url.Host {
	case "m.tv.sohu.com":
		data := sohuMVid.FindAllSubmatch([]byte(s.url.String()), 1)
		if len(data) == 0 {
			return nil, errors.New("未能在url中匹配到vid")
		}
		vid = string(data[0][1])
	case "tv.sohu.com":
		data := sohuTVid.FindAllSubmatch([]byte(s.url.String()), 1)
		if len(data) == 0 {
			return nil, errors.New("未能在url中匹配到vid")
		}
		paths, err := base64.StdEncoding.DecodeString(string(data[0][1]))
		if err != nil {
			return nil, err
		}
		vidPath := strings.Split(string(paths), "/")
		switch len(vidPath) {
		case 2: // 电视剧???
			return nil, errors.New("暂时不支持电影以及电视剧.")
		default:
			vid = strings.Split(strings.Split(string(paths), "/")[2], ".")[0]
		}
	}
	resp, err := s.client.Get(fmt.Sprintf(sohuVideoApi, vid), nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	sohuResp := new(sohuResponse)
	if err := json.Unmarshal(c, sohuResp); err != nil {
		return nil, err
	}
	if len(sohuResp.Data.ClipsURL) == 0 {
		return nil, errors.New("搜狐API未能返回视频下载url")
	}
	if sohuSuVid.Match([]byte(sohuResp.Data.Su[0])) {
		resp, err = s.client.Get(sohuNewVideoApi+sohuResp.Data.Su[0], nil)
		if err != nil {
			return nil, err
		}
		c, _ = ioutil.ReadAll(resp.Body)
		sohuNewResp := new(sohuNewResponse)
		if err := json.Unmarshal(c, sohuNewResp); err != nil {
			return nil, err
		}
		if len(sohuNewResp.Servers) == 0 || sohuNewResp.Servers[0].URL == "" {
			return nil, errors.New("搜狐API未能返回视频下载url")
		}
		return []download.Info{{Title: sohuResp.Data.TvName, Url: sohuNewResp.Servers[0].URL, Type: "mp4"}}, nil
	} else if sohuResp.Data.ClipsURL[0] != "" {
		return []download.Info{{Title: sohuResp.Data.TvName, Url: sohuResp.Data.ClipsURL[0], Type: "mp4"}}, nil
	}
	return nil, errors.New("搜狐API未能返回视频下载url")
}
