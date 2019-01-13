package extractors

import (
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
	"io/ioutil"
	url2 "net/url"
	"regexp"
	"strings"
)

type krcom struct {
	url    *url2.URL
	client *download.Http
}

func (k *krcom) Init(url string) error {
	k.url, _ = url2.Parse(url)
	k.client = download.NewHttp(nil)
	return nil
}

var (
	krcomVideo = regexp.MustCompile(`fluency=(.*?)\\"`)
	krcomTitle = regexp.MustCompile(`<title>(.*?)</title>`)
)

func (k *krcom) GetDownloadHeaders() map[string]string { return nil }
func (k *krcom) GetFileInfo() ([]download.Info, error) {
	data := make([]download.Info, 0)
	resp, err := k.client.Get(k.url.String(), nil)
	if err != nil {
		return data, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	for _, v := range krcomVideo.FindAllString(string(c), 1) {
		if v != "" {
			v = strings.Replace(v, `fluency=`, "", 1)
			v = strings.Replace(v, `\"`, "", 1)
			v, err := url2.PathUnescape(v)
			if err != nil {
				return data, err
			}
			var title string
			for _, t := range krcomTitle.FindAllString(string(c), 1) {
				if t != "" {
					title = strings.Replace(t, "<title>", "", 1)
					title = strings.Replace(title, "</title>", "", 1)
				}
			}
			if strings.HasPrefix(v, "//") {
				v = "http:" + v
			}
			data = append(data, download.Info{Title: title, Url: v})
			return data, nil
		}
	}
	return data, errors.New("酷然获取视频信息失败，规则变动。")
}
