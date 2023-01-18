package utils

import "strings"

func AcctInfo(acct string) (username, host string) {
	_acct := acct
	if strings.Contains(_acct, "acct:") {
		_acct = strings.TrimPrefix(_acct, "acct:")
	}
	if _acct[0] == '@' {
		_acct = _acct[1:]
	}
	if !strings.Contains(_acct, "@") {
		username = _acct
	} else {
		arr := strings.Split(_acct, "@")
		username = arr[0]
		host = arr[1]
	}
	return
}
