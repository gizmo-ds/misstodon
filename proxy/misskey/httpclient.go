package misskey

import (
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/go-resty/resty/v2"
)

var client = resty.New().
	SetHeader("User-Agent", "misstodon/"+global.AppVersion)
