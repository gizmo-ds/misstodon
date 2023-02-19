package misskey

import (
	"os"
	"testing"

	"github.com/gizmo-ds/misstodon/models"
	"github.com/stretchr/testify/assert"
)

func TestTimelinePublic(t *testing.T) {
	server := os.Getenv("TEST_SERVER")
	list, err := TimelinePublic(
		server, "",
		models.TimelinePublicTypeLocal, false,
		30, "", "")
	assert.NoError(t, err)
	t.Log(len(list))
	assert.NotEmpty(t, list)
}
