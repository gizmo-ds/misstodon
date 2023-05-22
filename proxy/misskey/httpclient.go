package misskey

import (
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func SetHeader(header, value string) {
	client.SetHeader(header, value)
}
