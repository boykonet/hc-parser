package http_requester

type IHTTPRequester interface {
	SetCookie(cookie string)
	SetHeaders(values map[string]string)
	SetRequestURI(uri string)
	GetBody() []byte
	ConfigureHTTP2Client() error
	Do() error
}