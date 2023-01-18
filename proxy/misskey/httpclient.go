package misskey

import (
	"github.com/go-resty/resty/v2"
)

var client *resty.Client

func SetClient(c *resty.Client) {
	client = c
}
