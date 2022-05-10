package model

import "time"

type AlarmNotice struct {
	Id              uint64    `gorm:"column:id;primary_key;autoIncrement"`
	AlarmRuleID     string    `gorm:"column:alarm_rule_id"`                            // 规则id
	ContractGroupID string    `gorm:"column:contract_group_id"`                        // 联系组id
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
}

func (*AlarmNotice) TableName() string {
	return "t_alarm_notice"
}
