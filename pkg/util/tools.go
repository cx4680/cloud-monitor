package util

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
)

func If(isTrue bool, a, b interface{}) interface{} {
	if isTrue {
		return a
	}
	return b
}

func GetDateDiff(diff int) string {
	nd := 1000 * 24 * 60 * 60
	nh := 1000 * 60 * 60
	nm := 1000 * 60
	ns := 1000
	// 计算差多少天
	day := diff / nd
	// 计算差多少小时
	hour := diff % nd / nh
	// 计算差多少分钟
	min := diff % nd % nh / nm
	// 计算差多少秒//输出结果
	sec := diff % nd % nh % nm / ns

	builder := strings.Builder{}

	if day > 0 {
		builder.WriteString(strconv.Itoa(day))
		builder.WriteString("天")
	}

	if hour > 0 {
		builder.WriteString(strconv.Itoa(hour))
		builder.WriteString("小时")
	}

	if min > 0 {
		builder.WriteString(strconv.Itoa(min))
		builder.WriteString("分钟")
	}

	if sec > 0 {
		builder.WriteString(strconv.Itoa(sec))
		builder.WriteString("秒")
	}
	str := builder.String()
	if strutil.IsNotBlank(str) {
		return str
	}
	return "0分钟"
}

func MD5(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	logger.Logger().Infof("md5 source=%s", data)
	h := md5.New()
	h.Write(data)
	ret := hex.EncodeToString(h.Sum(nil))
	logger.Logger().Infof("md5 result=%s", ret)
	return ret, nil
}

func SecToTime(seconds int) string {
	days := seconds / (3600 * 24)
	hour := (seconds - days*(3600*24)) / 3600
	minute := (seconds - hour*3600) / 60
	second := seconds - hour*3600 - minute*60
	sb := strings.Builder{}
	if days > 0 {
		sb.WriteString(strconv.Itoa(days) + "d")
	}
	if hour > 0 {
		sb.WriteString(strconv.Itoa(hour) + "h")
	}
	if minute > 0 {
		sb.WriteString(strconv.Itoa(minute) + "m")
	}
	if second > 0 {
		sb.WriteString(strconv.Itoa(second) + "s")
	}
	return sb.String()
}
