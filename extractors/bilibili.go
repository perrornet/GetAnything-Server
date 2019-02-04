package extractors

import (
	"errors"
	"fmt"
	"github.com/PerrorOne/GetAnything-Server/download"
	"github.com/PerrorOne/GetAnything-Server/utils"
	"io/ioutil"
	url2 "net/url"
	"regexp"
	"strings"
)

var (
	bilibiliCid   = regexp.MustCompile(`","cid":(\d+),`)
	bilibiliUrl   = regexp.MustCompile(`\<url\>\<!\[CDATA\[(.*?)\]\]`)
	bilibiliTitle = regexp.MustCompile(`"title":"(.*?)",`)
)

type bilibili struct {
	url     *url2.URL
	client  *download.Http
	headers map[string]string
	SEC1    string
	ApiUrl  string
}

func (b *bilibili) Init(url string) error {
	b.url, _ = url2.Parse(url)
	b.headers = map[string]string{"referer": url, "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36"}
	b.client = download.NewHttp(b.headers, true)
	b.SEC1 = "94aba54af9065f71de72f5508f1cd42e"
	b.ApiUrl = "http://interface.bilibili.com/v2/playurl?"
	return nil
}
func (b *bilibili) GetDownloadHeaders() map[string]string {
	return b.headers
}
func (b *bilibili) GetFileInfo() ([]download.Info, error) {
	switch b.url.Host {
	case "www.bilibili.com":
		resp, err := b.client.Get(b.url.String(), nil)
		if err != nil {
			return nil, err
		}
		c, _ := ioutil.ReadAll(resp.Body)
		cids := bilibiliCid.FindAll(c, 1)
		if len(cids) == 0 {
			return nil, errors.New("未能再当前页面寻找到cid")
		}
		titles := bilibiliTitle.FindAll(c, 1)
		var title string
		if len(titles) > 0 {
			fmt.Println()
			title = strings.Replace(string(titles[0]), `"title":"`, "", 1)
			title = strings.Replace(title, `",`, "", 1)
		}
		cid := strings.Replace(string(cids[0]), `","cid":`, "", 1)
		cid = strings.Replace(cid, `,`, "", 1)
		params := fmt.Sprintf("appkey=84956560bc028eb7&cid=%s&otype=xml&qn=%d&quality=%d&type=", cid, 15, 15)
		url := b.ApiUrl + params + "&sign=" + utils.NewMd5(params+b.SEC1).Encrypt()
		resp, err = b.client.Get(url, nil)
		if err != nil {
			return nil, err
		}
		c, _ = ioutil.ReadAll(resp.Body)
		urls := bilibiliUrl.FindAll(c, 1)
		if len(urls) == 0 {
			return nil, errors.New("未能找到相关的视频下载URL")
		}
		downloadUrl := strings.Replace(string(urls[0]), `<url><![CDATA[`, "", 1)
		downloadUrl = strings.Replace(downloadUrl, `]]`, "", 1)
		data := make([]download.Info, 0)
		data = append(data, download.Info{Title: title, Url: downloadUrl, Type: "flv"})
		return data, nil
	}
	return nil, nil
	//chksum = hashlib.md5(bytes(params_str+self.SEC1, 'utf8')).hexdigest()
	//api_url = self.api_url + params_str + '&sign=' + chksum
	// http://interface.bilibili.com/v2/playurl?appkey=84956560bc028eb7&cid=74472811&otype=xml&qn=15&quality=15&type=&sign=a6ff389a7c169890d531d060f88b2eb8
	// http://bangumi.bilibili.com/player/web_api/playurl?appkey=84956560bc028eb7&cid=74472811&otype=xml&qn=15&quality=15&type=&sign=a6ff389a7c169890d531d060f88b2eb8
}
