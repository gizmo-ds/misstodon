package httpclient

import (
	"io"
	"net/http"
)

type (
	Client interface {
		SetHeader(header, value string)
		R() Request
	}
	Request interface {
		SetBaseURL(url string) Request
		SetBody(body any) Request
		SetQueryParam(param string, value string) Request
		SetFormData(data map[string]string) Request
		SetDoNotParseResponse(parse bool) Request
		SetMultipartField(param string, fileName string, contentType string, reader io.Reader) Request
		SetResult(res any) Request
		Get(url string) (Response, error)
		Post(url string) (Response, error)
		Patch(url string) (Response, error)
	}
	Response interface {
		RawBody() io.ReadCloser
		Body() []byte
		String() string
		StatusCode() int
		Header() http.Header
	}
)
