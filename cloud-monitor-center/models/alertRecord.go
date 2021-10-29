package models

import "time"

// 告警记录
type AlertRecord struct {
	ID           string    `gorm:"column:id;primary_key"`
	Status       string    `gorm:"column:status;NOT NULL"` // 告警状态 firing resolved
	TenantID     string    `gorm:"column:tenant_id"`
	RuleID       string    `gorm:"column:rule_id;NOT NULL"` // 规则id
	RuleName     string    `gorm:"column:rule_name;NOT NULL"`
	MonitorType  string    `gorm:"column:monitor_type"`         // 监控类型
	SourceType   string    `gorm:"column:source_type;NOT NULL"` // 资源类型
	SourceID     string    `gorm:"column:source_id;NOT NULL"`   // 资源id
	Summary      string    `gorm:"column:summary"`
	CurrentValue string    `gorm:"column:current_value;NOT NULL"` // 当前值
	StartTime    time.Time `gorm:"column:start_time;NOT NULL"`    // 告警开始时间
	EndTime      time.Time `gorm:"column:end_time"`               // 告警结束时间
	TargetValue  string    `gorm:"column:target_value;NOT NULL"`  // 规则定义阈值
	Expression   string    `gorm:"column:expression;NOT NULL"`    // 计算公式
	Duration     string    `gorm:"column:duration;NOT NULL"`      // 持续时间
	Level        int       `gorm:"column:level;NOT NULL"`         // 告警级别 紧急 重要 次要 提醒
	NoticeStatus string    `gorm:"column:notice_status;NOT NULL"` // 消息发送状态, error 发送失败 success 成功
	AlarmKey     string    `gorm:"column:alarm_key;NOT NULL"`     // 告警项
	ContactInfo  string    `gorm:"column:contact_info"`           // 联系方式
	Region       string    `gorm:"column:region"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime;default:time.now()"`
	UpdateTime   time.Time `gorm:"column:update_time;autoCreateTime;default:time.now()"`
}

func (*AlertRecord) TableName() string {
	return "t_alert_record"
}
