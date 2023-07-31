package misskey_test

import (
	"testing"

	"github.com/gizmo-ds/misstodon/internal/utils"
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

func TestAccountMute(t *testing.T) {
	if _, ok := utils.StrEvaluation(testServer, testUserID, testToken); !ok {
		t.Skip("TEST_SERVER and TEST_USER_ID and TEST_TOKEN are required")
	}

	err := misskey.AccountMute(testServer, testToken, testUserID, 10*60)
	assert.NoError(t, err)

	account, err := misskey.AccountsGet(testServer, testToken, testUserID)
	assert.NoError(t, err)
	assert.Equal(t, true, *account.Limited)
}

func TestAccountUnmute(t *testing.T) {
	if _, ok := utils.StrEvaluation(testServer, testUserID, testToken); !ok {
		t.Skip("TEST_SERVER and TEST_USER_ID and TEST_TOKEN are required")
	}

	err := misskey.AccountUnmute(testServer, testToken, testUserID)
	assert.NoError(t, err)

	account, err := misskey.AccountsGet(testServer, testToken, testUserID)
	assert.NoError(t, err)
	assert.Equal(t, false, *account.Limited)
}
