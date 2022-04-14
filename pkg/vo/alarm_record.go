package vo

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
)

type AlarmRecordPageVO struct {
	Id           string          `json:"id"`
	Status       string          `json:"status"`
	TenantId     string          `json:"tenantId"`
	BizId        string          `json:"bizId"`
	RuleId       string          `json:"ruleId"`
	RuleName     string          `json:"ruleName"`
	MonitorType  string          `json:"monitorType"`
	SourceType   string          `json:"sourceType"`
	SourceId     string          `json:"sourceId"`
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

type AlarmRecordDetailVO struct {
	Id           string          `json:"id"`
	Status       string          `json:"status"`
	TenantId     string          `json:"tenantId"`
	BizId        string          `json:"bizId"`
	RuleId       string          `json:"ruleId"`
	RuleName     string          `json:"ruleName"`
	MonitorType  string          `json:"monitorType"`
	SourceType   string          `json:"sourceType"`
	SourceId     string          `json:"sourceId"`
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

type RecordNumHistory struct {
	Number  int    `json:"number"`
	DayTime string `json:"dayTime"`
}
