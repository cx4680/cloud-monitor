package models

type AlertRecord struct {
	Id           string `orm:"id" json:"id"`
	Status       string `orm:"status" json:"status"` // 告警状态 firing resolved
	TenantId     string `orm:"tenant_id" json:"tenant_id"`
	RuleId       string `orm:"rule_id" json:"rule_id"` // 规则id
	RuleName     string `orm:"rule_name" json:"rule_name"`
	MonitorType  string `orm:"monitor_type" json:"monitor_type"` // 监控类型
	SourceType   string `orm:"source_type" json:"source_type"`   // 资源类型
	SourceId     string `orm:"source_id" json:"source_id"`       // 资源id
	Summary      string `orm:"summary" json:"summary"`
	CurrentValue string `orm:"current_value" json:"current_value"`                // 当前值
	StartTime    string `orm:"start_time" json:"start_time"`                      // 告警开始时间
	EndTime      string `orm:"end_time" json:"end_time"`                          // 告警结束时间
	TargetValue  string `orm:"target_value" json:"target_value"`                  // 规则定义阈值
	Expression   string `orm:"expression" json:"expression" gorm:"size:500"`      // 计算公式
	Duration     string `orm:"duration" json:"duration"`                          // 持续时间
	Level        int    `orm:"level" json:"level"`                                // 告警级别 紧急 重要 次要 提醒
	NoticeStatus string `orm:"notice_status" json:"notice_status"`                // 消息发送状态, error 发送失败 success 成功
	AlarmKey     string `orm:"alarm_key" json:"alarm_key"`                        // 告警项
	ContactInfo  string `orm:"contact_info" json:"contact_info" gorm:"size:1000"` // 联系方式
	Region       string `orm:"region" json:"region"`
	CreateTime   string `orm:"create_time" json:"create_time"`
	UpdateTime   string `orm:"update_time" json:"update_time"`
}

func (*AlertRecord) TableName() string {
	return "t_alert_record"
}
