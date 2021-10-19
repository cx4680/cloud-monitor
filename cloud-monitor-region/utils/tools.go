package utils

func If(isTrue bool, a, b interface{}) interface{} {
	if isTrue {
		return a
	}
	return b
}
