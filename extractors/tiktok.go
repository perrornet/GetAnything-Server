package extractors

import (
	"GetAnything-Server/download"
	error2 "GetAnything-Server/error"
	"io/ioutil"
	"regexp"
	"strings"
)

type tiktok struct {
	content string
}

var (
	tiktokTitle = regexp.MustCompile(`<p class="desc">(.*?)</p>`)
	tiktokVideo = regexp.MustCompile(`playAddr: "(.*?)"`)
)

func (d *tiktok) GetFileFormUrl(url string) (string, error) {
	h := download.NewHttp(nil)
	resp, err := h.Get(url, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	c, _ := ioutil.ReadAll(resp.Body)
	d.content = string(c)
	for _, t := range tiktokVideo.FindAllString(d.content, 1) {
		if t != "" {
			return strings.Split(t, `"`)[1], nil
		}
	}
	return "", error2.NotFoundHandlerError
}
func (d *tiktok) GetFileInfo() *download.Info {
	for _, c := range tiktokTitle.FindAllString(d.content, 1) {
		if c != "" {
			c = strings.Replace(c, `<p class="desc">`, "", 1)
			c = strings.Replace(c, "</p>", "", 1)
			return &download.Info{Title: c}
		}
	}
	return &download.Info{}
}
