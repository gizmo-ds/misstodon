package utils

import "strings"

func Contains[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func StrEvaluation(str ...string) (v string, ok bool) {
	for _, s := range str {
		if s != "" {
			return s, true
		}
	}
	return
}

func Unique[T comparable](list []T) []T {
	var result []T
	t := make(map[T]struct{})
	for _, e := range list {
		t[e] = struct{}{}
	}
	for e := range t {
		result = append(result, e)
	}
	return result
}

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

func GetMentions(text string) []string {
	var result []string
	for _, s := range strings.Split(text, " ") {
		if strings.HasPrefix(s, "@") {
			result = append(result, s)
		}
	}
	return result
}
