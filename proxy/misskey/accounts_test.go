package misskey_test

import (
	"testing"

	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	if testServer == "" || testAcct == "" {
		t.Skip("TEST_SERVER and TEST_ACCT are required")
	}
	info, err := misskey.AccountsLookup(testServer, testAcct)
	assert.NoError(t, err)
	assert.Equal(t, testAcct, info.Acct)
}
