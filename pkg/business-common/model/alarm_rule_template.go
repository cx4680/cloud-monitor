package model

type AlarmRuleTemplate struct {
	Id             uint64 `gorm:"id" json:"id"`
	BizId          string `gorm:"biz_id" json:"biz_id"`
	MonitorType    string `gorm:"monitor_type" json:"monitor_type"`     // 监控类型
	ProductBizId   string `gorm:"product_biz_id" json:"product_biz_id"` // 产品名称
	Name           string `gorm:"name" json:"name"`                     // 规则名称
	MetricCode     string `gorm:"metric_code" json:"metric_code"`
	SilencesTime   string `gorm:"silences_time" json:"silences_time"`     // 冷却周期
	EffectiveStart string `gorm:"effective_start" json:"effective_start"` // 监控时间段-开始时间
	EffectiveEnd   string `gorm:"effective_end" json:"effective_end"`     // 监控时间段-结束时间
	Level          uint8  `gorm:"level" json:"level"`                     // 报警级别 紧急1 重要2 次要3 提醒 4
	CreateTime     string `gorm:"create_time" json:"create_time"`         // 创建时间
	Type           uint8  `gorm:"type" json:"type"`                       // 1:单指标, 2:多指标
	Combination    uint8  `gorm:"combination" json:"combination"`         // 逻辑组合：1:或, 2:且
	Period         int    `gorm:"period" json:"period"`
	Times          int    `gorm:"times" json:"times"`
}

func (*AlarmRuleTemplate) TableName() string {
	return "t_alarm_rule_template"
}
