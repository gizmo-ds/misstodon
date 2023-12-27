package misskey_test

import (
	"testing"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/provider/misskey"
	"github.com/stretchr/testify/assert"
)

func TestTimelinePublic(t *testing.T) {
	if testServer == "" {
		t.Skip("TEST_SERVER is required")
	}
	ctx := misstodon.ContextWithValues(testServer, testToken)
	list, err := misskey.TimelinePublic(
		ctx,
		models.TimelinePublicTypeLocal, false,
		30, "", "")
	assert.NoError(t, err)
	t.Log(len(list))
}
