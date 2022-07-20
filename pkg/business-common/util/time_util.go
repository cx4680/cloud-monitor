package util

import "time"

const (
	FullTimeFmt = "2006-01-02 15:04:05"
	ZoneTimeFmt = "2006-01-02T15:04:05Z"
	DayTimeFmt  = "2006-01-02"
)

// GetMonthStartEnd 获取指定时间所在月的开始 结束时间
func GetMonthStartEnd(t time.Time) (time.Time, time.Time) {
	monthStartDay := t.AddDate(0, 0, -t.Day()+1)
	monthStartTime := time.Date(monthStartDay.Year(), monthStartDay.Month(), monthStartDay.Day(), 0, 0, 0, 0, t.Location())
	monthEndDay := monthStartTime.AddDate(0, 1, -1)
	monthEndTime := time.Date(monthEndDay.Year(), monthEndDay.Month(), monthEndDay.Day(), 23, 59, 59, 0, t.Location())
	return monthStartTime, monthEndTime
}

func GetNow() time.Time {
	return time.Now()
}

func GetNowStr() string {
	return time.Now().Format(FullTimeFmt)
}

func TimeToStr(t time.Time, fmt string) string {
	return t.Format(fmt)
}

func StrToTime(fmt, str string) (t time.Time) {
	t, _ = time.Parse(fmt, str)
	return
}

func TimeParseForZone(s string) time.Time {
	parse, _ := time.Parse(ZoneTimeFmt, s)
	return time.Unix(parse.Unix(), 0)
}

func TimeToFullTimeFmtStr(t time.Time) string {
	return t.Format(FullTimeFmt)
}

func TimestampToFullTimeFmtStr(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(FullTimeFmt)
}

func TimestampToDayTimeFmtStr(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(DayTimeFmt)
}
