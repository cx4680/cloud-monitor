package models

type AlarmNotice struct {
	AlarmRuleId     string `orm:"alarm_rule_id" json:"alarm_rule_id"`         // 规则id
	ContractGroupId string `orm:"contract_group_id" json:"contract_group_id"` // 联系组id
	CreateTime      string `orm:"create_time" json:"create_time"`
}

func (*AlarmNotice) TableName() string {
	return "t_alarm_notice"
}
