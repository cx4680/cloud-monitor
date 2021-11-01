package models

import "time"

type AlarmInstance struct {
	AlarmRuleID  string    `gorm:"column:alarm_rule_id"`
	InstanceID   string    `gorm:"column:instance_id"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime;default:time.now()"` // 创建时间
	RegionCode   string    `gorm:"column:region_code"`
	ZoneCode     string    `gorm:"column:zone_code"`
	Ip           string    `gorm:"column:ip"`
	RegionName   string    `gorm:"column:region_name"`
	ZoneName     string    `gorm:"column:zone_name"`
	InstanceName string    `gorm:"column:instance_name"`
	TenantID     string    `gorm:"column:tenant_id"` // 租户id
}

func (*AlarmInstance) TableName() string {
	return "t_alarm_instance"
}
