package extractors

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PerrorOne/GetAnything-Server/download"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	toutiaoHeaders = map[string]string{"user-agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"}
	toutiaoVid     = regexp.MustCompile(`videoId\s*:\s*'([^']+)'`)
	toutiaoTitle   = regexp.MustCompile(`title: '(\S+)',`)
)

type toutiao struct {
	client  *download.Http
	url     string
	content string
}

type toutiaoResp struct {
	Data struct {
		VideoList struct {
			Video1 struct {
				MainURL string `json:"main_url"`
			} `json:"video_1"`
		} `json:"video_list"`
	} `json:"data"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Total   int    `json:"total"`
}

//代码来自https://github.com/soimort/you-get/blob/develop/src/you_get/extractors/toutiao.py
func signVideoUrl(vid string) string {
	rr := rand.New(rand.NewSource(time.Now().Unix()))
	r := strconv.FormatFloat(rr.Float64(), 'f', 10, 64)
	r = r[2:]
	url := "http://i.snssdk.com/video/urls/v/1/toutiao/mp4/" + vid
	n := strings.Replace(url, "http://i.snssdk.com", "", 1) + "?r=" + r
	c := crc32.ChecksumIEEE([]byte(n))
	s := c >> 0
	return url + fmt.Sprintf("?r=%s&s=%d", r, s)
}

func (t *toutiao) Init(url string) error {
	t.url = url
	t.client = download.NewHttp(nil)
	return nil
}

func (t *toutiao) GetDownloadHeaders() map[string]string { return nil }

func (t *toutiao) GetFileInfo() ([]download.Info, error) {
	data := make([]download.Info, 0)
	vid, err := t.ixigua(t.url)
	if err != nil {
		return data, err
	}
	url, err := t.GetUrlFromVid(vid)
	if err != nil {
		return data, err
	}
	c := toutiaoTitle.FindAllString(t.content, 1)
	for _, i := range c {
		if i != "" {
			i = strings.Replace(i, "title: '", "", 1)
			data = append(data, download.Info{Title: strings.Replace(i, "',", "", 1), Url: url})
			return data, nil
		}
	}
	return data, errors.New("头条系视频接口变化。")
}

func (t *toutiao) GetUrlFromVid(vid string) (string, error) {
	url := signVideoUrl(vid)
	fmt.Println(url)
	resp, err := t.client.Get(url, nil)
	if err != nil {
		return "", err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	data := &toutiaoResp{}
	if err := json.Unmarshal(c, data); err != nil {
		return "", err
	}
	if data.Code != 0 {
		return "", errors.New(data.Message)
	}
	fmt.Println(data.Data.VideoList.Video1.MainURL)
	decoded, err := base64.StdEncoding.DecodeString(data.Data.VideoList.Video1.MainURL)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// https://www.ixigua.com/
func (t *toutiao) ixigua(url string) (string, error) {
	t.client = download.NewHttp(toutiaoHeaders)
	resp, err := t.client.Get(url, nil)
	if err != nil {
		return "", err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	t.content = string(c)
	vids := toutiaoVid.FindAllString(t.content, 1)
	if len(vids) == 0 {
		return "", errors.New("未能获取到正确的vid")
	}
	vids = strings.Split(vids[0], "'")
	if len(vids) < 2 {
		return "", errors.New("未能匹配到正确的vid")
	}
	return vids[1], nil
}
