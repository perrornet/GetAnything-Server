package extractors

import (
	"GetAnything-Server/download"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

type tiktok struct {
	content string
	url     string
}

var (
	tiktokTitle = regexp.MustCompile(`<p class="desc">(.*?)</p>`)
	tiktokVideo = regexp.MustCompile(`playAddr: "(.*?)"`)
)

func (d *tiktok) Init(url string) error {
	d.url = url
	return nil
}

func (d *tiktok) GetDownloadHeaders() map[string]string { return nil }

func (d *tiktok) GetFileInfo() ([]download.Info, error) {
	h := download.NewHttp(nil)
	data := make([]download.Info, 0)
	resp, err := h.Get(d.url, nil)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	c, _ := ioutil.ReadAll(resp.Body)
	var title string
	for _, c := range tiktokTitle.FindAllString(string(c), 1) {
		if c != "" {
			c = strings.Replace(c, `<p class="desc">`, "", 1)
			title = strings.Replace(c, "</p>", "", 1)
		}
	}
	for _, t := range tiktokVideo.FindAllString(d.content, 1) {
		if t != "" {
			data = append(data, download.Info{Url: strings.Split(t, `"`)[1], Title: title})
			return data, nil
		} else {
			return data, errors.New("未能正确获取抖音视频下载URL")
		}
	}
	return data, errors.New("未能正确获取抖音视频下载URL")
}
