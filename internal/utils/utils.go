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
