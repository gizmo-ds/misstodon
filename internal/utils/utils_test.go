package utils_test

import (
	"testing"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestAcctInfo(t *testing.T) {
	username, host := utils.AcctInfo("@gizmo_ds@liuli.lol")
	assert.Equal(t, username, "gizmo_ds")
	assert.Equal(t, host, "liuli.lol")

	username, host = utils.AcctInfo("@banana")
	assert.Equal(t, username, "banana")
	assert.Equal(t, host, "")

	username, host = utils.AcctInfo("user@misskey.io")
	assert.Equal(t, username, "user")
	assert.Equal(t, host, "misskey.io")
}

func TestGetMentions(t *testing.T) {
	mentions := utils.GetMentions("Hello @gizmo_ds@liuli.lol")
	assert.Equal(t, len(mentions), 1)
	assert.Equal(t, mentions[0], "@gizmo_ds@liuli.lol")

	mentions = utils.GetMentions("@user@misskey.io")
	assert.Equal(t, len(mentions), 1)
	assert.Equal(t, mentions[0], "@user@misskey.io")

	mentions = utils.GetMentions("@banana")
	assert.Equal(t, len(mentions), 1)
	assert.Equal(t, mentions[0], "@banana")
}
