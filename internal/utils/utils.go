package utils

import (
	"sort"
	"strconv"
	"strings"
)

// Contains returns true if the list contains the item, false otherwise.
func Contains[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

// StrEvaluation evaluates a list of strings and returns the first non-empty
// string, or an empty string if no non-empty strings are found.
func StrEvaluation(str ...string) (v string, ok bool) {
	for _, s := range str {
		if s != "" {
			return s, true
		}
	}
	return
}

// Unique returns a new list containing only the unique elements of list.
// The order of the elements is preserved.
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

// AcctInfo splits an account string into a username and host.
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

// SliceIfNull returns the given slice if it is not nil, or an empty slice if it is nil.
func SliceIfNull[T any](slice []T) []T {
	if slice == nil {
		return []T{}
	}
	return slice
}

type accountField struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	VerifiedAt *string
}

// GetFieldsAttributes converts a map of fields to a slice of accountFields
// The map of fields is expected to have keys in the form fields_attributes[<index>][<tag>]
// Where <index> is an integer and <tag> is one of "name" or "value".
// The order of the accountFields in the returned slice is determined by the order of the <index> values.
func GetFieldsAttributes(values map[string][]string) (fields []accountField) {
	var m = make(map[int]*accountField)
	var keys []int
	for k, v := range values {
		ok, index, tag := func(field string) (ok bool, index int, tag string) {
			if len(field) < (17+3+6) ||
				field[:17] != "fields_attributes" ||
				field[17] != '[' ||
				field[len(field)-1] != ']' {
				return
			}
			field = field[18 : len(field)-1]
			if !strings.Contains(field, "][") {
				return
			}
			parts := strings.Split(field, "][")
			if len(parts) != 2 || (parts[0] == "" || parts[1] == "") {
				return
			}
			var err error
			index, err = strconv.Atoi(parts[0])
			if err != nil {
				return
			}
			ok = true
			tag = parts[1]
			return
		}(k)
		if !ok {
			continue
		}
		if _, e := m[index]; !e {
			m[index] = &accountField{}
		}
		switch tag {
		case "name":
			m[index].Name = v[0]
		case "value":
			m[index].Value = v[0]
		}
	}
	for i, f := range m {
		if f.Name == "" {
			continue
		}
		keys = append(keys, i)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fields = append(fields, *m[k])
	}
	return
}

func NumRangeLimit[T int | int64](i, min, max T) T {
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}
