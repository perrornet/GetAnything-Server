package extractors

import (
	"GetAnything-Server/download"
	"errors"
	"io/ioutil"
	url2 "net/url"
	"regexp"
	"strings"
)

type kuaishou struct {
	url    *url2.URL
	client *download.Http
}

var (
	kuaishouVideo   = regexp.MustCompile(`type="video/mp4" src=".*?\.mp4"`)
	kuaishouHeaders = map[string]string{"User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.81 Mobile Safari/537.36"}
	kuaishouTitle   = regexp.MustCompile(`<meta itemprop=name content=.*?>`)
)

func (k *kuaishou) Init(url string) error {
	k.url, _ = url2.Parse(url)
	k.client = download.NewHttp(kuaishouHeaders)
	return nil
}
func (k *kuaishou) GetDownloadHeaders() map[string]string { return nil }
func (k *kuaishou) GetFileInfo() ([]download.Info, error) {
	data := make([]download.Info, 0)
	resp, err := k.client.Get(k.url.String(), kuaishouHeaders)
	if err != nil {
		return data, err
	}
	var title string
	c, _ := ioutil.ReadAll(resp.Body)
	titles := kuaishouTitle.FindAll(c, 1)
	if len(titles) > 0 {
		title = strings.Replace(string(titles[0]), `<meta itemprop=name content=`, "", 1)
		title = strings.Replace(title, `>`, "", 1)
	}
	for _, t := range kuaishouVideo.FindAllString(string(c), 1) {
		if t != "" {
			t = strings.Replace(t, `type="video/mp4" src="`, "", 1)
			t = strings.Replace(t, `"`, "", 1)
			data = append(data, download.Info{Url: t, Title: title, Type: "mp4"})
			return data, nil
		}
	}
	return data, errors.New("快手视频url匹配失败。")
}
