package misskey

import (
	"github.com/gizmo-ds/misstodon/pkg/httpclient"
)

var client = httpclient.NewRestyClient()

func SetHeader(header, value string) {
	client.SetHeader(header, value)
}
