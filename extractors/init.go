package extractors

import (
	"github.com/PerrorOne/GetAnything-Server/download"
	error2 "github.com/PerrorOne/GetAnything-Server/error"
	url2 "net/url"
)

// m3u8格式视频返回格式需要将.ts文件的下载url按照顺序放入[]download.Info
// 并在type中表明m3u8
func newDownload(host string) download.Download {
	switch host {
	case "v.douyin.com", "www.iesdouyin.com":
		return &tiktok{}
	case "v.douyu.com", "vmobile.douyu.com":
		return &douyuTV{}
	case "www.ixigua.com", "www.365yg.com", "m.toutiaoimg.cn", "m.365yg.com":
		return &toutiao{}
	case "weibo.com", "m.weibo.cn":
		return &weibo{}
	case "www.meipai.com":
		return &meipai{}
	case "krcom.cn":
		return &krcom{}
	case "m.gifshow.com", "live.kuaishou.com", "www.gifshow.com":
		return &kuaishou{}
	case "n.miaopai.com":
		return &miaopai{}
	case "www.zhihu.com", "zhuanlan.zhihu.com":
		return &zhihu{}
	case "www.bilibili.com":
		return &bilibili{}
	case "h5.weishi.qq.com":
		return &weishi{}
	case "v.huya.com":
		return &huya{}
	default:
		return nil
	}
	return nil
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
