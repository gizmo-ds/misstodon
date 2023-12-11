package httpclient

import (
	"io"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/go-resty/resty/v2"
)

func NewRestyClient() Client {
	return &RestyClient{client: resty.New()}
}

type (
	RestyClient struct {
		client *resty.Client
	}
	RestyRequest struct {
		baseURL string
		req     *resty.Request
	}
)

func (c RestyClient) SetHeader(header, value string) {
	c.client.SetHeader(header, value)
}

func (c RestyClient) R() Request {
	return &RestyRequest{req: c.client.R()}
}

func (r RestyRequest) SetBaseURL(url string) Request {
	r.baseURL = url
	return r
}

func (r RestyRequest) SetBody(body any) Request {
	r.req = r.req.SetBody(body)
	return r
}

func (r RestyRequest) SetQueryParam(param, value string) Request {
	r.req = r.req.SetQueryParam(param, value)
	return r
}

func (r RestyRequest) SetFormData(data map[string]string) Request {
	r.req = r.req.SetFormData(data)
	return r
}

func (r RestyRequest) SetDoNotParseResponse(parse bool) Request {
	r.req = r.req.SetDoNotParseResponse(parse)
	return r
}

func (r RestyRequest) SetMultipartField(param, fileName, contentType string, reader io.Reader) Request {
	r.req = r.req.SetMultipartField(param, fileName, contentType, reader)
	return r
}

func (r RestyRequest) SetResult(res any) Request {
	r.req = r.req.SetResult(res)
	return r
}

func (r RestyRequest) Get(url string) (Response, error) {
	u := url
	if r.baseURL != "" {
		u = utils.JoinURL(r.baseURL, url)
	}
	return r.req.Get(u)
}

func (r RestyRequest) Post(url string) (Response, error) {
	u := url
	if r.baseURL != "" {
		u = utils.JoinURL(r.baseURL, url)
	}
	return r.req.Post(u)
}

func (r RestyRequest) Patch(url string) (Response, error) {
	u := url
	if r.baseURL != "" {
		u = utils.JoinURL(r.baseURL, url)
	}
	return r.req.Patch(u)
}

func (r RestyRequest) Delete(url string) (Response, error) {
	u := url
	if r.baseURL != "" {
		u = utils.JoinURL(r.baseURL, url)
	}
	return r.req.Delete(u)
}
