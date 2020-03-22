package extractors

import (
	"GetAnything-Server/download"
	"GetAnything-Server/logger"
	"encoding/json"
	"errors"
	"fmt"
	logger2 "github.com/apsdehal/go-logger"
	"io/ioutil"
	url1 "net/url"
	"regexp"
	"strings"
)

var (
	zhihuTitle       = regexp.MustCompile(`data-react-helmet="true">(.*?)</title>`)
	zhihuVideoList   = regexp.MustCompile(`<a class="video-box" href="\S+video/(\d+)"`)
	zhihuFakeHeaders = map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:60.0) Gecko/20100101 Firefox/60.0",
	}
	log = logger.NewLogger("GetAnything", logger2.InfoLevel)
)

type zhihuResponse struct {
	Playlist struct {
		Ld struct {
			PlayURL string `json:"play_url"`
		} `json:"ld"`
		Hd struct {
			PlayURL string `json:"play_url"`
		} `json:"hd"`
		Sd struct {
			PlayURL string `json:"play_url"`
		} `json:"sd"`
	} `json:"playlist"`
}

type zhihuQuestion struct {
	Data []struct {
		ID int `json:"id"`
	} `json:"data"`
	Paging struct {
		IsEnd bool   `json:"is_end"`
		Next  string `json:"next"`
	} `json:"paging"`
}

type zhihu struct {
	url    *url1.URL
	client *download.Http
}

func (z *zhihu) anser(url string, offset int) ([]download.Info, error) {
	if offset == 0 {
		offset++
	}
	u, _ := url1.Parse(url)
	paths := strings.Split(u.Path, "/")
	if len(paths) < 2 {
		return nil, errors.New("不支持该url")
	}
	resp, err := z.client.Get(url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	c, _ := ioutil.ReadAll(resp.Body)
	var title string
	t := zhihuTitle.FindAll(c, 1)
	if len(t) >= 1 {
		title = strings.Replace(string(t[0]), `data-react-helmet="true">`, "", 1)
		title = strings.Replace(title, `</title>`, "", 1)
	}
	videos := zhihuVideoList.FindAll(c, -1)
	if len(videos) <= 0 {
		return nil, errors.New("未能找到可用的视频ID")
	}
	data := make([]download.Info, 0)
	for index, i := range videos {
		u := strings.Replace(string(i), `<a class="video-box" href="`, "", 1)
		u = strings.Replace(string(i), `"`, "", 1)
		urls := strings.Split(u, "/")
		url := strings.Replace("https://lens.zhihu.com/api/videos/"+urls[len(urls)-1], `"`, "", 1)
		resp, err := z.client.Get(url, nil)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		c, _ := ioutil.ReadAll(resp.Body)
		r := &zhihuResponse{}
		if err = json.Unmarshal(c, r); err != nil {
			log.Error(err.Error())
			continue
		}
		if r.Playlist.Hd.PlayURL != "" {
			data = append(data, download.Info{Title: fmt.Sprintf("%s_%d_%d", title, offset, index), Url: r.Playlist.Hd.PlayURL})
		} else if r.Playlist.Ld.PlayURL != "" {
			data = append(data, download.Info{Title: fmt.Sprintf("%s_%d_%d", title, offset, index), Url: r.Playlist.Ld.PlayURL})
		} else if r.Playlist.Sd.PlayURL != "" {
			data = append(data, download.Info{Title: fmt.Sprintf("%s_%d_%d", title, offset, index), Url: r.Playlist.Sd.PlayURL})
		} else {
			continue
		}
	}
	if len(data) != 0 {
		return data, nil
	}
	return nil, errors.New("未能在该回答中寻找到视频。")
}

func (z *zhihu) zvideo() ([]download.Info, error) {
	resp, err := z.client.Get(z.url.String(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	c, _ := ioutil.ReadAll(resp.Body)
	var title string
	t := zhihuTitle.FindAll(c, 1)
	if len(t) >= 1 {
		title = strings.Replace(string(t[0]), `data-react-helmet="true">`, "", 1)
		title = strings.Replace(title, `</title>`, "", 1)
	}
	id := regexp.MustCompile(`\{"videoId":"(.*?)",`)
	vids := id.FindAllStringSubmatch(string(c), 1)
	if len(vids) < 1 {
		paths := strings.Split(z.url.Path, "/")
		return []download.Info{{Title: title, Url: fmt.Sprintf("https://video.zhihu.com/video/%s", paths[len(paths)-1])}}, nil
	}
	resp1, err := z.client.Get(fmt.Sprintf("https://lens.zhihu.com/api/v4/videos/%s", vids[0][1]), nil)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()
	var r zhihuResponse
	if _, err := download.GetJson(resp1.Body, &r); err != nil {
		return nil, err
	}
	data := make([]download.Info, 0)
	if r.Playlist.Hd.PlayURL != "" {
		data = append(data, download.Info{Title: title, Url: r.Playlist.Hd.PlayURL})
	} else if r.Playlist.Ld.PlayURL != "" {
		data = append(data, download.Info{Title: title, Url: r.Playlist.Ld.PlayURL})
	} else if r.Playlist.Sd.PlayURL != "" {
		data = append(data, download.Info{Title: title, Url: r.Playlist.Sd.PlayURL})
	} else {
		return nil, errors.New("视频不存在或者规则改变")
	}
	return data, nil
}

func (z *zhihu) Init(url string) error {
	z.url, _ = url1.Parse(url)
	z.client = download.NewHttp(zhihuFakeHeaders, true)
	return nil
}
func (z *zhihu) GetDownloadHeaders() map[string]string {
	return nil
}
func (z *zhihu) GetFileInfo() ([]download.Info, error) {
	switch z.url.Host {
	case "www.zhihu.com":
		switch {
		case strings.Contains(z.url.Path, "zvideo"):
			return z.zvideo()
		case strings.Contains(z.url.Path, "answer"):
			return z.anser(z.url.String(), 0)
		}
		return nil, errors.New("目前禁止使用， 当一个问题下有多个回答可能会消耗过多的时间去请求，导致无谓的网络IO")
	case "zhuanlan.zhihu.com":
		return z.anser(z.url.String(), 0)
	}
	return nil, nil
}
