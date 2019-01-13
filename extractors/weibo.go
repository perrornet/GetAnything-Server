package extractors

import (
	"errors"
	"fmt"
	"github.com/PerrorOne/GetAnything-Server/download"
	"io/ioutil"
	url2 "net/url"
	"regexp"
	"strings"
)

type weibo struct {
	content string
	url     string
	client  *download.Http
}

var (
	//videoInfoUrl   = "https://m.weibo.cn/s/video/object?object_id=%s&mid=%s"
	weiboTitle1 = regexp.MustCompile(`"title": "(.*?)",`)
	weiboVideo1 = regexp.MustCompile(`"mp4_hd_mp4": "(.*?)",`)
	weiboUrl    = "https://m.weibo.cn/status/%s?type=comment&jumpfrom=weibocom"
)

func (w *weibo) Init(url string) error {
	w.url = url
	w.client = download.NewHttp(nil)
	return nil
}

func (w *weibo) GetDownloadHeaders() map[string]string { return nil }

func (w *weibo) GetFileInfo() ([]download.Info, error) {
	data := make([]download.Info, 0)
	u, _ := url2.Parse(w.url)
	switch u.Host {
	case "weibo.com", "m.weibo.cn":
		ids := strings.Split(u.Path, "/")
		if len(ids) < 3 {
			return data, errors.New("Url错误未找到该微博ID")
		}
		id := ids[2]
		client := download.NewHttp(nil)
		resp, err := client.Get(fmt.Sprintf(weiboUrl, id), nil)
		if err != nil {
			return data, err
		}
		c, _ := ioutil.ReadAll(resp.Body)
		w.content = string(c)
		var title string
		for _, t := range weiboTitle1.FindAllString(w.content, 1) {
			if t != "" {
				title = strings.Replace(t, `"title": "`, "", 1)
				title = strings.Replace(title, `",`, "", 1)
				break
			}
		}
		for _, v := range weiboVideo1.FindAllString(w.content, 1) {
			if v != "" {
				v = strings.Replace(v, `"mp4_hd_mp4": "`, "", 1)
				v = strings.Replace(v, `",`, "", 1)
				data = append(data, download.Info{Url: v, Title: title})
				return data, nil
			}
			return data, errors.New("未能正确匹配到该微博下的视频下载URL")
		}
	}
	return data, errors.New("当前类型暂未支持！")
}
