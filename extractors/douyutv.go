package extractors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PerrorOne/GetAnything-Server/download"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
)

type douyuTV struct {
	content string
	url     string
}

var (
	douyuVideoInfoUrl = "http://vmobile.douyu.com/video/getInfo?vid="
	douyuTitle        = regexp.MustCompile("<h1>(.+?)</h1>")
	douyuHeaders      = map[string]string{"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"}
)

type douyuResp struct {
	Error int         `json:"error"`
	Data  interface{} `json:"data"`
}

func (d *douyuTV) Init(url string) error {
	d.url = url
	return nil
}

func (d *douyuTV) GetDownloadHeaders() map[string]string { return nil }

func (d *douyuTV) GetFileInfo() ([]download.Info, error) {
	if !strings.Contains(d.url, "show/") {
		return nil, errors.New("未能提供正确的UR，URL中必须包含'show/'")
	}
	d.url = strings.Replace(d.url, "vmobile.douyu.com", "v.douyu.com", 1)
	t := strings.Split(d.url, "/")
	if len(t) < 5 {
		return nil, errors.New("未能提供正确的URL，示例UR:" +
			"https://vmobile.douyu.com/show/8pa9v5pL91KWVrqA?source=qq&medium=and&type=vd")
	}
	vId := strings.Split(t[4], "?")[0]
	client := download.NewHttp(douyuHeaders, true)
	resp, err := client.Get(d.url, nil)
	if err != nil {
		return nil, err
	}
	c, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	d.content = string(c)
	resp, err = client.Get(douyuVideoInfoUrl+vId, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData := &douyuResp{}
	c, _ = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(c, respData); err != nil {
		return nil, err
	}
	if respData.Error != 0 {
		return nil, errors.New(respData.Data.(string))
	}
	var title string
	titles := douyuTitle.FindAllString(d.content, 1)
	for _, t := range titles {
		title = strings.Replace(t, "<h1>", "", 1)
		title = strings.Replace(title, "</h1>", "", 1)
		break
	}
	data := make([]download.Info, 0)
	if v, ok := respData.Data.(map[string]interface{}); ok {
		if v1, ok := v["video_url"]; ok {
			resp, err = client.Get(v1.(string), nil)
			if err != nil {
				return nil, err
			}
			m3u8Url, _ := url.Parse(v1.(string))
			downloadUrl := m3u8Url.Host + strings.Join(strings.Split(m3u8Url.Path, "/")[:3], "/")
			c, _ = ioutil.ReadAll(resp.Body)
			for _, u := range bytes.Split(c, []byte("\n")) {
				if bytes.HasPrefix(u, []byte("#")) || bytes.Equal(u, []byte("")) {
					continue
				}
				data = append(data, download.Info{Title: title, Url: url.QueryEscape(fmt.Sprintf("http://%s/%s", downloadUrl, string(u))), Type: "m3u8"})
			}
			return data, nil
		} else {
			return data, errors.New("请求斗鱼视频信息获取失败，data.data中没有'video_url'键")
		}
	}
	return nil, errors.New("斗鱼获取视频信息接口类型变化。")
}
