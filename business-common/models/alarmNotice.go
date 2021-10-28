package models

import "time"

type AlarmNotice struct {
	AlarmRuleId     string    `orm:"alarm_rule_id" json:"alarm_rule_id"`                                       // 规则id
	ContractGroupId string    `orm:"contract_group_id" json:"contract_group_id"`                               // 联系组id
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime;default:time.now()" json:"create_time"` // 创建时间
}

func (*AlarmNotice) TableName() string {
	return "t_alarm_notice"
}
