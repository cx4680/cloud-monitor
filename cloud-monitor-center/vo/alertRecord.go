package vo

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
)

type AlertRecordPageVO struct {
	ID           string
	Status       string
	TenantID     string
	RuleID       string
	RuleName     string
	MonitorType  string
	SourceType   string
	SourceID     string
	Summary      string
	CurrentValue string
	StartTime    global.JsonTime
	EndTime      global.JsonTime
	TargetValue  string
	Expression   string
	Duration     string
	Level        int
	NoticeStatus string
	AlarmKey     string
	ContactInfo  string
	Region       string
	CreateTime   global.JsonTime
	UpdateTime   global.JsonTime
}

type AlertRecordDetailVO struct {
	ID           string
	Status       string
	TenantID     string
	RuleID       string
	RuleName     string
	MonitorType  string
	SourceType   string
	SourceID     string
	Summary      string
	CurrentValue string
	StartTime    global.JsonTime
	EndTime      global.JsonTime
	TargetValue  string
	Expression   string
	Duration     string
	Level        int
	NoticeStatus string
	AlarmKey     string
	ContactInfo  string
	Region       string
	CreateTime   global.JsonTime
	UpdateTime   global.JsonTime
}

type RecordNumHistory struct {
	Number  int    `json:"number"`
	DayTime string `json:"dayTime"`
}
