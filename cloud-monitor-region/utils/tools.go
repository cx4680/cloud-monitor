package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"
)

func If(isTrue bool, a, b interface{}) interface{} {
	if isTrue {
		return a
	}
	return b
}

func MD5(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func SecToTime(seconds int) string {
	days := seconds / (3600 * 24)
	hour := (seconds - days*(3600*24)) / 3600
	minute := (seconds - hour*3600) / 60
	second := seconds - hour*3600 - minute*60
	sb := strings.Builder{}
	if days > 0 {
		sb.WriteString(string(days) + "d")
	}
	if hour > 0 {
		sb.WriteString(string(hour) + "h")
	}
	if minute > 0 {
		sb.WriteString(string(minute) + "m")
	}
	if second > 0 {
		sb.WriteString(string(second) + "s")
	}
	return sb.String()
}
