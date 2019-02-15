package extractors

import (
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
	"github.com/PerrorOne/GetAnything-Server/utils"
	"io/ioutil"
	"regexp"
)

var (
	kgqqPlayVideoUrl = regexp.MustCompile(`"playurl_video":"(?s:(.*?))"`)
	kgqqPlayM4aUrl   = regexp.MustCompile(`"playurl":"(?s:(.*?))"`)
)

type kgqq struct {
	url    string
	client *download.Http
}

// 最先调用此方法， 该方法建议只初始化一些参数
func (k *kgqq) Init(url string) error {
	k.url = url
	k.client = download.NewHttp(nil)
	return nil
}

// 获取下载文件时所需的headers, 如果没有返回nil
func (k *kgqq) GetDownloadHeaders() map[string]string { return nil }

// 获取文件的下载url
func (k *kgqq) GetFileInfo() ([]download.Info, error) {
	resp, err := k.client.Get(k.url, nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	var url, _type string
	videos := kgqqPlayVideoUrl.FindAllSubmatch(c, 1)
	if string(videos[0][1]) == "" {
		m4as := kgqqPlayM4aUrl.FindAllSubmatch(c, 1)
		if len(m4as) == 0 {
			return nil, errors.New("未能在页面中寻找到可下载的URL")
		}
		url = string(m4as[0][1])
		_type = "m4a"
	} else {
		url = string(videos[0][1])
		_type = "mp4"
	}
	m := utils.NewMd5(c)
	return []download.Info{{Title: m.Encrypt(), Url: url, Type: _type}}, nil
}
