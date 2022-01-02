package model

import (
	"time"
)

type AlertRecord struct {
	Id           string    `gorm:"id;primary_key" json:"id"`
	Status       string    `gorm:"status" json:"status"` // 告警状态 firing resolved
	TenantId     string    `gorm:"tenant_id" json:"tenant_id"`
	RuleId       string    `gorm:"rule_id" json:"rule_id"` // 规则id
	RuleName     string    `gorm:"rule_name" json:"rule_name"`
	MonitorType  string    `gorm:"monitor_type" json:"monitor_type"` // 监控类型
	SourceType   string    `gorm:"source_type" json:"source_type"`   // 资源类型
	SourceId     string    `gorm:"source_id" json:"source_id"`       // 资源id
	Summary      string    `gorm:"summary" json:"summary"`
	CurrentValue string    `gorm:"current_value" json:"current_value"`                // 当前值
	StartTime    string    `gorm:"start_time" json:"start_time"`                      // 告警开始时间
	EndTime      string    `gorm:"end_time" json:"end_time"`                          // 告警结束时间
	TargetValue  string    `gorm:"target_value" json:"target_value"`                  // 规则定义阈值
	Expression   string    `gorm:"expression" json:"expression" gorm:"size:500"`      // 计算公式
	Duration     string    `gorm:"duration" json:"duration"`                          // 持续时间
	Level        uint8     `gorm:"level" json:"level"`                                // 告警级别 紧急 重要 次要 提醒
	NoticeStatus string    `gorm:"notice_status" json:"notice_status"`                // 消息发送状态, error 发送失败 success 成功
	AlarmKey     string    `gorm:"alarm_key" json:"alarm_key"`                        // 告警项
	ContactInfo  string    `gorm:"contact_info" json:"contact_info" gorm:"size:1000"` // 联系方式
	Region       string    `gorm:"region" json:"region"`
	CreateTime   time.Time `gorm:"create_time" json:"create_time"`
	UpdateTime   time.Time `gorm:"update_time" json:"update_time"`
}

func (*AlertRecord) TableName() string {
	return "t_alert_record"
}
