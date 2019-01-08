package extractors

import (
	"GetAnything-Server/download"
	error2 "GetAnything-Server/error"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ixiguaHeaders = map[string]string{"user-agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"}
	ixiguaVid     = regexp.MustCompile(`videoId\s*:\s*'([^']+)'`)
	ixiguaTitle   = regexp.MustCompile(`title: '(\S+)',`)
)

type toutiao struct {
	content string
	client  *download.Http
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
	fmt.Println(r)
	url := "http://i.snssdk.com/video/urls/v/1/toutiao/mp4/" + vid
	n := strings.Replace(url, "http://i.snssdk.com", "", 1) + "?r=" + r
	c := crc32.ChecksumIEEE([]byte(n))
	s := c >> 0
	return url + fmt.Sprintf("?r=%s&s=%d", r, s)
}

func (t *toutiao) GetUrlFromvid(vid string) (string, error) {
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
	t.client = download.NewHttp(ixiguaHeaders)
	resp, err := t.client.Get(url, nil)
	if err != nil {
		return "", err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	t.content = string(c)
	vids := ixiguaVid.FindAllString(t.content, 1)
	if len(vids) == 0 {
		return "", error2.UrlError
	}
	vids = strings.Split(vids[0], "'")
	if len(vids) < 2 {
		return "", error2.UrlError
	}
	return vids[1], nil
}

func (t *toutiao) GetFileFormUrl(url string) (string, error) {
	vid, err := t.ixigua(url)
	if err != nil {
		return "", err
	}
	return t.GetUrlFromvid(vid)
}

func (d *toutiao) GetFileInfo() *download.Info { return &download.Info{} }
