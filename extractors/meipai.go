package extractors

import (
	"encoding/base64"
	"errors"
	"github.com/PerrorOne/GetAnything-Server/download"
	"github.com/antchfx/xpath"
	"github.com/antchfx/xquery/html"
	"io/ioutil"
	url2 "net/url"
	"regexp"
	"strconv"
	"strings"
)

type meipai struct { // 解密代码来自https://blog.csdn.net/Crazy__Hope/article/details/84995399，原py脚本
	url    *url2.URL
	client *download.Http
}

var (
	meipaiDataVideo = regexp.MustCompile(`data-video="(.*?)">`)
	meipaiTitle     = xpath.MustCompile(`/html/body/div[3]/div[2]/div[2]/h1[1]/text()`)
)

func (m *meipai) Init(url string) error {
	m.url, _ = url2.Parse(url)
	m.client = download.NewHttp(nil)
	return nil
}

func (m *meipai) getHex(dataVideo string) (string, string) {
	hex := dataVideo[:4]
	reHex := make([]byte, 4, 4)
	for index, v := range hex {
		reHex[3-index] = uint8(v)
	}
	hex = string(reHex)
	return dataVideo[4:], hex
}

func (m *meipai) getDec(hex string) ([]string, []string) {
	k, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return []string{}, []string{}
	}
	kStr := strconv.Itoa(int(k))
	t := []byte(kStr[:2])
	t1 := []byte(kStr[2:])
	return []string{string(t[0]), string(t[1])}, []string{string(t1[0]), string(t1[1])}
}

func (m meipai) substr(s string, s2 []string) string {
	k, err := strconv.Atoi(s2[0])
	if err != nil {
		return ""
	}
	c := s[:k]
	addNum, err := strconv.Atoi(s2[1])
	if err != nil {
		return ""
	}
	d := s[k : k+addNum]
	temp := strings.Replace(s[k:], d, "", 1)
	return c + temp
}

func (m *meipai) getPos(s string, s1 []string) []string {
	s1_0, err := strconv.Atoi(s1[0])
	if err != nil {
		return []string{}
	}
	s1_1, err := strconv.Atoi(s1[1])
	if err != nil {
		return []string{}
	}
	s1[0] = strconv.Itoa(len(s) - s1_0 - s1_1)
	return s1
}

func (m *meipai) GetDownloadHeaders() map[string]string { return nil }
func (m *meipai) GetFileInfo() ([]download.Info, error) {
	data := make([]download.Info, 0)
	switch m.url.Host {
	case "www.meipai.com": // PC
		resp, err := m.client.Get(m.url.String(), nil)
		if err != nil {
			return data, err
		}
		c, _ := ioutil.ReadAll(resp.Body)
		var dataVideo string
		for _, v := range meipaiDataVideo.FindAllString(string(c), 1) {
			if v != "" {
				dataVideo = strings.Replace(v, `data-video="`, "", 1)
				dataVideo = strings.Replace(dataVideo, `">`, "", 1)
			}
		}
		if dataVideo == "" {
			return data, errors.New("未能匹配到美拍视频url")
		}
		s, hex := m.getHex(dataVideo)
		pre, tail := m.getDec(hex)
		dd := m.substr(s, pre)
		v := m.getPos(dd, tail)
		if len(v) == 0 {
			return data, errors.New("计算视频真实URL出错，data-video=" + dataVideo)
		}
		videoUrl, err := base64.StdEncoding.DecodeString(m.substr(dd, v))
		if err != nil {
			return data, err
		}
		var title string
		doc, err := htmlquery.Parse(strings.NewReader(string(c)))
		if err == nil {
			iter := meipaiTitle.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
			for iter.MoveNext() {
				title = iter.Current().Value()
				title = strings.Replace(title, "\n", "", -1)
				title = strings.Replace(title, " ", "", -1)
			}
		}
		data = append(data, download.Info{Url: string(videoUrl), Title: title})
		return data, nil
	}
	return data, nil
}
