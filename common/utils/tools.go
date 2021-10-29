package utils

import "strconv"

func Int64ToInt(data int64) int {
	string := strconv.FormatInt(data, 10)
	result, _ := strconv.Atoi(string)
	return result
}
