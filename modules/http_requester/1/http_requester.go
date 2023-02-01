package http_requester

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"
	"vbtor/modules/http_requester"
)

type requester struct {
	Request *http.Request
	Response *http.Response

	Host *http.Client
}

func NewHTTP1Requester(method, url string, body []byte) http_requester.IHTTPRequester {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		
	}
	return &requester{
		Request: request,
		Host: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (r *requester) SetCookie(cookie string) {
	r.Request.Header.Add("cookie", cookie)
}

func (r *requester) SetHeaders(values map[string]string) {
	for key, value := range values {
		r.Request.Header.Add(key, value)
	}
}

func (r *requester) SetRequestURI(uri string) {
	r.Request.URL = &url.URL{Path: uri}
}

func (r *requester) Do() error {
	resp, err := r.Host.Do(r.Request)
	r.Response = resp
	return err
}

func (r *requester) GetBody() []byte {
	buff := bytes.Buffer{}
	io.Copy(&buff, r.Response.Body)
	return buff.Bytes()
}

func (r *requester) ConfigureHTTP2Client() error {
	panic("doesn't implemet")
}
