package models

type UsePrometheusId struct {
	TenantId         string `orm:"tenant_id" json:"tenant_id"`
	PrometheusRuleId string `orm:"prometheus_rule_id" json:"prometheus_rule_id"`
	CreateTime       string `orm:"create_time" json:"create_time"`
	CreateUser       string `orm:"create_user" json:"create_user"`
	Deleted          int    `orm:"deleted" json:"deleted"` // 1（已删除） 0未删除
}

func (*UsePrometheusId) TableName() string {
	return "t_user_prometheus_id"
}
