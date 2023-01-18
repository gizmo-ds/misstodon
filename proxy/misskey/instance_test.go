package misskey_test

import (
	"os"
	"testing"

	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	server := os.Getenv("TEST_SERVER")
	info, err := misskey.Instance(server, "development")
	assert.NoError(t, err)
	assert.Equal(t, server, info.Uri)
}
