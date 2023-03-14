package misskey_test

import (
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/go-resty/resty/v2"
)

func init() {
	misskey.SetClient(resty.New())
}
