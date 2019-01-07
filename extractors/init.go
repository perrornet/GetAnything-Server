package extractors

import (
	"GetAnything-Server/download"
	"errors"
	url2 "net/url"
)

var Site = map[string]download.Download{
	"v.douyin.com": &tiktok{},
}
func Match(url string)(download.Download, error){
	u, err := url2.Parse(url)
	if err != nil{
		return nil, err
	}
	if v, ok := Site[u.Host]; ok{
		return v, nil
	}
	return nil,errors.New("Not found handler function")
}
