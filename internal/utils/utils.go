package utils

func Contains[T string](list []T, item T) bool {
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

func Unique[T string](list []T) []T {
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
