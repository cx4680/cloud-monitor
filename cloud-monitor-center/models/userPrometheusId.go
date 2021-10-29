package models

import "time"

type UserPrometheusID struct {
	TenantID         string    `gorm:"column:tenant_id"`
	PrometheusRuleID string    `gorm:"column:prometheus_rule_id;NOT NULL"`
	CreateTime       time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	CreateUser       string    `gorm:"column:create_user"`
	Deleted          int       `gorm:"column:deleted;default:0"` // 1（已删除） 0未删除
}

func (*UserPrometheusID) TableName() string {
	return "t_user_prometheus_id"
}
