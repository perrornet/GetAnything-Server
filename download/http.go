package download

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// 文件的信息
type Info struct {
	// 文件下载的url
	// m3u8类型详见douyutv.go
	Url string `json:"url"`
	// 文件的标题， 文件将以标题命名
	Title string `json:"title"`

	// 该文件的类型,常规类型下文件将以该参数为后缀, m3u8类型文件标明m3u8
	// 如:mp4, flv, mov...

	Type string `json:"type"`
}

type Download interface {
	// 最先调用此方法， 该方法建议只初始化一些参数
	Init(url string) error
	// 获取下载文件时所需的headers, 如果没有返回nil
	GetDownloadHeaders() map[string]string
	// 获取文件的下载url
	GetFileInfo() ([]Info, error)
}

// 生成post所需的body数据
func postData(body map[string]string) io.Reader {
	if body == nil {
		return nil
	}
	data := make([]string, len(body), len(body))
	var index int
	for k, v := range body {
		data[index] = fmt.Sprintf("%s=%s", k, v)
		index++
	}
	return strings.NewReader(strings.Join(data, "&"))
}

type Http struct { // only request api, download video use tcp
	isSession bool
	client    *http.Client
	headers   map[string]string
	cookie    map[string]*http.Cookie
}

func NewHttp(headers map[string]string, isSession ...bool) *Http {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}, Timeout: 2 * time.Second}
	if len(isSession) == 0 {
		return &Http{client: client, headers: headers, cookie: map[string]*http.Cookie{}}
	}
	return &Http{client: client, isSession: true, headers: headers, cookie: map[string]*http.Cookie{}}
}

func (h *Http) request(method, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	} else if h.headers != nil {
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
	}

	for _, v := range h.cookie {
		req.AddCookie(v)
	}
	return req, nil
}

func (h *Http) do(method, url string, headers map[string]string, data map[string]string) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}
	req, err := h.request(method, url, postData(data), headers)
	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if req == nil {
		return nil, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	if h.isSession {
		for _, c := range req.Cookies() {
			h.cookie[c.Name] = c
		}
	}
	return resp, nil
}

// if headers != nil, cover Http.headers
func (h *Http) Get(url string, headers map[string]string) (*http.Response, error) {
	return h.do("GET", url, headers, nil)
}

func (h *Http) Head(url string, headers map[string]string) (*http.Response, error) {
	return h.do("HEAD", url, headers, nil)
}

func (h *Http) Post(url string, headers map[string]string, data map[string]string) (*http.Response, error) {
	return h.do("POST", url, headers, data)
}
