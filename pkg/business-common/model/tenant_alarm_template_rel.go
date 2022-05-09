package model

type TenantAlarmTemplateRel struct {
	Id            uint64 `gorm:"id" json:"id"`
	TenantId      string `gorm:"tenant_id" json:"tenant_id"`
	TemplateBizId string `gorm:"template_biz_id" json:"template_biz_id"`
	CreateTime    string `gorm:"create_time" json:"create_time"`
}

func (*TenantAlarmTemplateRel) TableName() string {
	return "t_tenant_alarm_template_rel"
}
