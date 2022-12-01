package http_requester

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/dgrr/http2"
)

type requester struct {
	Request *fasthttp.Request
	Response *fasthttp.Response

	Host *fasthttp.HostClient
}

func NewHTTPRequester(addrSSL string) IHTTPRequester {
	return &requester{
		Request: fasthttp.AcquireRequest(),
		Response: fasthttp.AcquireResponse(),
		Host: &fasthttp.HostClient{
			Addr: addrSSL,
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
	r.Request.SetRequestURI(uri)
}

func (r *requester) Do() error {
	return r.Host.Do(r.Request, r.Response)
}

func (r *requester) GetBody() []byte {
	return r.Response.Body()
}

func (r *requester) ConfigureHTTP2Client() error {
	err := http2.ConfigureClient(r.Host, http2.ClientOpts{})
	if err != nil {
		return fmt.Errorf("%s doesn't support http/2: %v", r.Host.Addr, err)
	}
	return nil
}