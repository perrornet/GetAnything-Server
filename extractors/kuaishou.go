package extractors

import (
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
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
	kuaishouVideo   = regexp.MustCompile(`src="http://.*?\.mp4"`)
	kuaishouHeaders = map[string]string{"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"}
)

func (k *kuaishou) Init(url string) error {
	k.url, _ = url2.Parse(url)
	k.client = download.NewHttp(kuaishouHeaders)
	return nil
}
func (k *kuaishou) GetDownloadHeaders() map[string]string { return nil }
func (k *kuaishou) GetFileInfo() ([]download.Info, error) {
	data := make([]download.Info, 0)
	resp, err := k.client.Get(k.url.String(), nil)
	if err != nil {
		return data, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	for _, t := range kuaishouVideo.FindAllString(string(c), 1) {
		if t != "" {
			t = strings.Replace(t, `src="`, "", 1)
			t = strings.Replace(t, `"`, "", 1)
			data = append(data, download.Info{Url: t})
			return data, nil
		}
	}
	return data, errors.New("快手视频url匹配失败。")
}
