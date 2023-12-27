package misskey

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAcctIsInvalid = errors.New("acct format is invalid")
	ErrRateLimit     = errors.New("rate limit")
)
