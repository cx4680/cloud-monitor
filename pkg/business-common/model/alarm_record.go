package model

import (
	"time"
)

type AlarmRecord struct {
	Id             uint64    `gorm:"id;primary_key;autoIncrement" json:"id"`
	BizId          string    `gorm:"column:biz_id;size=500" json:"bizId"`
	RequestId      string    `gorm:"request_id;size=50" json:"requestId"`
	Status         string    `gorm:"status" json:"status"` // 告警状态 firing resolved
	TenantId       string    `gorm:"tenant_id" json:"tenantId"`
	RuleId         string    `gorm:"rule_id" json:"ruleId"` // 规则id
	RuleName       string    `gorm:"rule_name" json:"ruleName"`
	RuleSourceType uint8     `gorm:"rule_source_type;default:1;force" json:"ruleSourceType"`
	RuleSourceId   string    `gorm:"rule_source_id" json:"ruleSourceId"`
	MonitorType    string    `gorm:"monitor_type" json:"monitorType"`   // 监控类型
	SourceType     string    `gorm:"source_type" json:"sourceType"`     // 资源类型
	SourceId       string    `gorm:"source_id" json:"sourceId"`         // 资源id
	CurrentValue   string    `gorm:"current_value" json:"currentValue"` // 当前值
	StartTime      time.Time `gorm:"start_time" json:"startTime"`       // 告警开始时间
	EndTime        time.Time `gorm:"end_time" json:"endTime"`           // 告警结束时间
	TargetValue    string    `gorm:"target_value" json:"targetValue"`   // 规则定义阈值
	Duration       string    `gorm:"duration" json:"duration"`          // 持续时间
	Level          uint8     `gorm:"level" json:"level"`                // 告警级别 紧急 重要 次要 提醒
	AlarmKey       string    `gorm:"alarm_key" json:"alarmKey"`         // 告警项
	Region         string    `gorm:"region" json:"region"`
	CreateTime     time.Time `gorm:"create_time" json:"createTime"`
	UpdateTime     time.Time `gorm:"update_time" json:"updateTime"`
}

func (*AlarmRecord) TableName() string {
	return "t_alarm_record"
}
