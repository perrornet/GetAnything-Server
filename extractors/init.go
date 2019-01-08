package extractors

import (
	"GetAnything-Server/download"
	error2 "GetAnything-Server/error"
	url2 "net/url"
)

var Site = map[string]download.Download{
	"v.douyin.com":      &tiktok{},
	"v.douyu.com":       &douyuTV{},
	"vmobile.douyu.com": &douyuTV{},
}

func Match(url string) (download.Download, error) {
	u, err := url2.Parse(url)
	if err != nil {
		return nil, error2.UrlError
	}
	if v, ok := Site[u.Host]; ok {
		return v, nil
	}
	return nil, error2.NotFoundHandlerError
}
