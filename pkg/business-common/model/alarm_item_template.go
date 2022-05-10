package model

type AlarmItemTemplate struct {
	Id                uint64     `gorm:"id" json:"id"`
	RuleTemplateBizId string     `gorm:"rule_template_biz_id" json:"rule_template_biz_id"`
	MetricCode        string     `gorm:"metric_code" json:"metric_code"`
	TriggerCondition  *Condition `gorm:"trigger_condition" json:"trigger_condition"` // 条件表达式
	Level             uint8      `gorm:"level" json:"level"`
	SilencesTime      string     `gorm:"silences_time" json:"silences_time"` // 告警间隔
}

func (*AlarmItemTemplate) TableName() string {
	return "t_alarm_item_template"
}
