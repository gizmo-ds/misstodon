package misskey_test

import (
	"testing"

	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/stretchr/testify/assert"
)

func TestTimelinePublic(t *testing.T) {
	if testServer == "" {
		t.Skip("TEST_SERVER is required")
	}
	list, err := misskey.TimelinePublic(
		testServer, testToken,
		models.TimelinePublicTypeLocal, false,
		30, "", "")
	assert.NoError(t, err)
	t.Log(len(list))
}
