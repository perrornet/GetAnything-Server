package extractors

import (
	"github.com/PerrorOne/GetAnything-Server/download"
	error2 "github.com/PerrorOne/GetAnything-Server/error"
	url2 "net/url"
)

func newDownload(host string) download.Download {
	switch host {
	case "v.douyin.com", "www.iesdouyin.com": // example url:http://v.douyin.com/Nxsx71/
		return &tiktok{}
	case "v.douyu.com", "vmobile.douyu.com": // example url:https://v.douyu.com/show/NbwE7ZBr8rB7n5Zz
		return &douyuTV{}
	case "www.ixigua.com", "www.365yg.com", "m.toutiaoimg.cn": // example url:http://www.365yg.com/a6642859345774117383/#mid=1616102707166216
		return &toutiao{}
	case "weibo.com", "m.weibo.cn": // example url:https://weibo.com/1739046981/HbbxNh0PO?type=comment
		return &weibo{}
	case "www.meipai.com": // example url:https://www.meipai.com/media/1051281695
		return &meipai{}
	case "krcom.cn": // example url:https://krcom.cn/6441132677/episodes/2358773:4326651149133458
		return &krcom{}
	case "m.gifshow.com", "live.kuaishou.com", "www.gifshow.com": // example url:http://www.gifshow.com/s/4wl4pk3y
		return &kuaishou{}
	case "n.miaopai.com": // example url:http://n.miaopai.com/media/pU-xTQ-A2keFNWqqwC3CVUwFjoywScpn.htm
		return &miaopai{}
	case "www.zhihu.com", "zhuanlan.zhihu.com": // example url:https://www.zhihu.com/question/282693696/answer/538355040
		return &zhihu{}
	case "www.bilibili.com":
		return &bilibili{}
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
