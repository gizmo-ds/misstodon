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

func TestGetFieldsAttributes(t *testing.T) {
	values := map[string][]string{
		"display_name":                {"cy"},
		"note":                        {"hello"},
		"fields_attributes[0][name]":  {"GitHub"},
		"fields_attributes[0][value]": {"https://github.com"},
		"fields_attributes[1][name]":  {"Twitter"},
		"fields_attributes[1][value]": {"https://twitter.com"},
		"fields_attributes[3][name]":  {"Google"},
		"fields_attributes[3][value]": {"https://google.com"},
		"fields_attributes[name]":     {"Google"},
		"fields_attributes[value]":    {"https://google.com"},
	}
	fields := utils.GetFieldsAttributes(values)
	assert.Equal(t, 3, len(fields))
	assert.Equal(t, "GitHub", fields[0].Name)
	assert.Equal(t, "https://github.com", fields[0].Value)
	assert.Equal(t, "Twitter", fields[1].Name)
	assert.Equal(t, "https://twitter.com", fields[1].Value)
	assert.Equal(t, "Google", fields[2].Name)
	assert.Equal(t, "https://google.com", fields[2].Value)
}
