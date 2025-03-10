package k8s

import "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"

type AlertRuleDTO struct {
	AlertRuleId    string       `json:"alertRuleId"`
	TenantId       string       `json:"tenantId"`
	Region         string       `json:"region"`         //所属区域
	Zone           string       `json:"zone"`           // 可用区
	SpecGroupsList []SpecGroups `json:"specGroupsList"` //告警规则详情
}

type SpecGroups struct {
	Name      string     `json:"name"`
	AlertList []AlertDTO `json:"alertList"`
}

type AlertDTO struct {
	RuleType     string                 `json:"ruleType"` // alert or record
	Alert        string                 `json:"alert"`
	Record       string                 `json:"record"`
	Expr         string                 `json:"expr"`
	ForTime      string                 `json:"for_time"` // 3m
	Summary      string                 `json:"summary"`
	Description  string                 `json:"description"` //告警详情
	Labels       map[string]interface{} `json:"labels"`
	SilencesTime string                 `json:"silencesTime"`
	SourceType   uint8                  `json:"source_type"`
}

type AlarmDescription struct {
	Expr            string
	ExprDetail      string
	Rule            model.AlarmRule
	RuleItems       []model.AlarmItem
	ContactGroupIds []string
	ResourceId      string
	ResourceGroupId string
}
