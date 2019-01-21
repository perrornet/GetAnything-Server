package extractors

import (
	"encoding/json"
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
	mWeiboUrl   = "https://m.weibo.cn/s/video/object?object_id=%s&mid=%s"
)

type weiboResp struct {
	Ok   int `json:"ok"`
	Data struct {
		ObjectType string `json:"object_type"`
		Object     struct {
			Summary string `json:"summary"`
			Stream  struct {
				URL string `json:"url"`
			} `json:"stream"`
		} `json:"object"`
	} `json:"data"`
}

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
	case "weibo.com":
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
	case "m.weibo.cn":
		u, _ := url2.Parse(w.url)
		objectId, mid := u.Query().Get("object_id"), u.Query().Get("blog_mid")
		if objectId == "" || mid == "" {
			return data, errors.New("未能在URL中寻找到object_id和blog_mid参数")
		}
		resp, err := w.client.Get(fmt.Sprintf(mWeiboUrl, objectId, mid), nil)
		if err != nil {
			return data, err
		}
		r := &weiboResp{}
		c, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(c, r); err != nil {
			return data, err
		}
		if r.Ok != 1 {
			return data, errors.New("微博接口返回失败， 未能找到对应的微博")
		}
		if r.Data.ObjectType != "video" {
			return data, errors.New("该条微博不包含视频。")
		}
		var title string
		t := strings.Split(r.Data.Object.Summary, "<span")
		if len(t) != 0 {
			title = t[0]
		}
		data = append(data, download.Info{Title: title, Url: r.Data.Object.Stream.URL})
		return data, nil
	}
	return data, errors.New("当前类型暂未支持！")
}
