package vo

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
)

type AlertRecordPageVO struct {
	ID           string          `json:"id"`
	Status       string          `json:"status"`
	TenantID     string          `json:"tenantId"`
	RuleID       string          `json:"ruleId"`
	RuleName     string          `json:"ruleName"`
	MonitorType  string          `json:"monitorType"`
	SourceType   string          `json:"source_type"`
	SourceID     string          `json:"sourceId"`
	Summary      string          `json:"summary"`
	CurrentValue string          `json:"currentValue"`
	StartTime    string          `json:"startTime"`
	EndTime      string          `json:"endTime"`
	TargetValue  string          `json:"targetValue"`
	Expression   string          `json:"expression"`
	Duration     string          `json:"duration"`
	Level        int             `json:"level"`
	NoticeStatus string          `json:"noticeStatus"`
	AlarmKey     string          `json:"alarmKey"`
	ContactInfo  string          `json:"contactInfo"`
	Region       string          `json:"region"`
	CreateTime   global.JsonTime `json:"createTime"`
	UpdateTime   global.JsonTime `json:"updateTime"`
}

type AlertRecordDetailVO struct {
	ID           string          `json:"id"`
	Status       string          `json:"status"`
	TenantID     string          `json:"tenantId"`
	RuleID       string          `json:"ruleId"`
	RuleName     string          `json:"ruleName"`
	MonitorType  string          `json:"monitorType"`
	SourceType   string          `json:"source_type"`
	SourceID     string          `json:"sourceId"`
	Summary      string          `json:"summary"`
	CurrentValue string          `json:"currentValue"`
	StartTime    global.JsonTime `json:"startTime"`
	EndTime      global.JsonTime `json:"endTime"`
	TargetValue  string          `json:"targetValue"`
	Expression   string          `json:"expression"`
	Duration     string          `json:"duration"`
	Level        int             `json:"level"`
	NoticeStatus string          `json:"noticeStatus"`
	AlarmKey     string          `json:"alarmKey"`
	ContactInfo  string          `json:"contactInfo"`
	Region       string          `json:"region"`
	CreateTime   global.JsonTime `json:"createTime"`
	UpdateTime   global.JsonTime `json:"updateTime"`
}

type RecordNumHistory struct {
	Number  int    `json:"number"`
	DayTime string `json:"dayTime"`
}
