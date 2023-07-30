package misskey_test

import (
	"testing"

	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	if testServer == "" {
		t.Skip("TEST_SERVER is required")
	}
	info, err := misskey.Instance(testServer, "development")
	assert.NoError(t, err)
	assert.Equal(t, testServer, info.Uri)
}
