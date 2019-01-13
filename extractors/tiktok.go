package extractors

import (
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
	"io/ioutil"
	"regexp"
	"strings"
)

type tiktok struct {
	content string
	url     string
}

var (
	tiktokTitle   = regexp.MustCompile(`<p class="desc">(.*?)</p>`)
	tiktokVideo   = regexp.MustCompile(`playAddr: "(.*?)"`)
	tiktokHeaders = map[string]string{"user-agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"}
)

func (d *tiktok) Init(url string) error {
	d.url = url
	return nil
}

func (d *tiktok) GetDownloadHeaders() map[string]string { return nil }

func (d *tiktok) GetFileInfo() ([]download.Info, error) {
	h := download.NewHttp(tiktokHeaders)
	data := make([]download.Info, 0)
	resp, err := h.Get(d.url, nil)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	c, _ := ioutil.ReadAll(resp.Body)
	d.content = string(c)
	var title string
	for _, c := range tiktokTitle.FindAllString(d.content, 1) {
		if c != "" {
			c = strings.Replace(c, `<p class="desc">`, "", 1)
			title = strings.Replace(c, "</p>", "", 1)
		}
	}
	for _, t := range tiktokVideo.FindAllString(d.content, 1) {
		if t != "" {
			data = append(data, download.Info{Url: strings.Split(t, `"`)[1], Title: title})
			return data, nil
		}
	}
	return data, errors.New("未能正确获取抖音视频下载URL")
}
