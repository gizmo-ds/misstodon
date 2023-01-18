package misskey_test

import (
	"os"
	"testing"

	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	server := os.Getenv("TEST_SERVER")
	acct := os.Getenv("TEST_ACCT")
	info, err := misskey.Lookup(server, acct)
	assert.NoError(t, err)
	assert.Equal(t, acct, info.Acct)
}
