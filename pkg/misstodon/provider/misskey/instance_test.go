package misskey_test

import (
	"testing"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey"
	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	if testServer == "" {
		t.Skip("TEST_SERVER is required")
	}
	ctx := misstodon.ContextWithValues(testServer, "")
	info, err := misskey.Instance(ctx, "development")
	assert.NoError(t, err)
	assert.Equal(t, testServer, info.Uri)
}
