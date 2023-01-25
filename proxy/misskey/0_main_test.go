package misskey_test

import (
	"os"
	"testing"

	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/go-resty/resty/v2"
)

func TestMain(m *testing.M) {
	misskey.SetClient(resty.New())
	os.Exit(m.Run())
}
