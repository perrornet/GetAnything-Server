package extractors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PerrorOne/GetAnything-Server/download"
	error2 "github.com/PerrorOne/GetAnything-Server/error"
	"io/ioutil"
	url2 "net/url"
	"strings"
)

var (
	nMiaopaiComvideoinfo = "https://n.miaopai.com/api/aj_media/info.json?smid=%s" // https://n.miaopai.com/xxx
)

type miaopai struct {
	url    *url2.URL
	client *download.Http
}

type nMiaopaiComResp struct {
	Code int `json:"code"`
	Data struct {
		Description string `json:"description"`
		MetaData    []struct {
			PlayUrls struct {
				N string `json:"n"`
			} `json:"play_urls"`
		} `json:"meta_data"`
	} `json:"data"`
	Key interface{} `json:"key"`
}

func (m *miaopai) nMiaopaiCom() ([]download.Info, error) {
	data := make([]download.Info, 0)
	paths := strings.Split(m.url.Path, "/")
	if len(paths) < 3 {
		return data, errors.New("url错误，示例url:" + "http://n.miaopai.com/media/pU-xTQ-A2keFNWqqwC3CVUwFjoywScpn.htm")
	}
	smid := strings.Split(paths[2], ".")[0]
	resp, err := m.client.Get(fmt.Sprintf(nMiaopaiComvideoinfo, smid), nil)
	if err != nil {
		return data, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	r := &nMiaopaiComResp{}
	if err := json.Unmarshal(c, r); err != nil {
		return data, err
	}

	if r.Code != 200 {
		return data, errors.New("请求视频信息错误，接口返回失败")
	}
	for _, v := range r.Data.MetaData {
		if v.PlayUrls.N != "" {
			data = append(data, download.Info{Title: r.Data.Description, Url: v.PlayUrls.N})
			return data, nil
		}
	}
	return data, errors.New("请求视频信息错误，接口返回数据中没有视频信息")
}

func (m *miaopai) Init(url string) error {
	m.url, _ = url2.Parse(url)
	m.client = download.NewHttp(nil)
	return nil
}
func (m *miaopai) GetDownloadHeaders() map[string]string { return nil }

func (m *miaopai) GetFileInfo() ([]download.Info, error) {
	switch m.url.Host {
	case "n.miaopai.com":
		return m.nMiaopaiCom()
	}
	return []download.Info{}, error2.NotFoundHandlerError
}
