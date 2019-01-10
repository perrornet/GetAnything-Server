package download

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"time"
)

type Info struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type Download interface {
	Init(url string) error
	GetDownloadHeaders() map[string]string
	GetFileInfo() ([]Info, error)
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

func (h *Http) do(method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}
	req, err := h.request(method, url, body, headers)
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
// if error return nil
func (h *Http) Get(url string, headers map[string]string) (*http.Response, error) {
	return h.do("GET", url, headers, nil)
}

func (h *Http) Post(url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	return h.do("POST", url, headers, body)
}
