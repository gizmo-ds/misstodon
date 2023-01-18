package misskey_test

import (
	"os"
	"testing"

	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")
	misskey.SetClient(resty.New())
	os.Exit(m.Run())
}
