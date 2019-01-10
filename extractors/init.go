package extractors

import (
	"GetAnything-Server/download"
	error2 "GetAnything-Server/error"
	url2 "net/url"
)

func newDownload(host string) download.Download {
	switch host {
	case "v.douyin.com":
		return &tiktok{}
	case "v.douyu.com", "vmobile.douyu.com":
		return &douyuTV{}
	case "www.ixigua.com", "www.365yg.com", "m.toutiaoimg.cn":
		return &toutiao{}
	case "weibo.com", "m.weibo.cn":
		return &weibo{}
	case "www.meipai.com":
		return &meipai{}
	case "krcom.cn":
		return &krcom{}
	default:
		return nil
	}
}

func Match(url string) (download.Download, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return nil, err
	}
	if d := newDownload(u.Host); d != nil {
		if err := d.Init(url); err != nil {
			return nil, err
		}
		return d, nil
	}
	return nil, error2.NotFoundHandlerError
}
