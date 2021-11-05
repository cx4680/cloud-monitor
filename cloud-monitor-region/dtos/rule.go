package dtos

import "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"

type RuleExpress struct {
	RuleId        string
	RuleName      string
	ProductType   string
	MonitorType   string
	Level         int
	NoticeChannel int
	RuleCondition *forms.RuleCondition
	GroupIds      []*forms.NoticeGroup
	InstanceList  []*forms.InstanceInfo
}

type RuleDescDTO struct {
	RuleName           string   `json:"ruleName "`
	Product            string   `json:"product"`
	InstanceInfo       string   `json:"instanceInfo"`
	MetricName         string   `json:"metricName "`
	ComparisonOperator string   `json:"comparisonOperator"`
	TargetValue        int      `json:"targetValue "`
	Time               int      `json:"time"`
	Period             int      `json:"period"`
	Unit               string   `json:"unit"`
	Express            string   `json:"express"`
	Level              int      `json:"level"`
	NotifyChannel      string   `json:"notifyChannel"`
	MonitorItem        string   `json:"monitorItem"`
	MonitorType        string   `json:"monitorType"`
	InstanceId         string   `json:"instanceId"`
	RuleId             string   `json:"ruleId"`
	TenantId           string   `json:"tenantId"`
	RegionName         string   `json:"regionName"`
	Statistic          string   `json:"statistic"`
	GroupList          []string `json:"groupList"`
}
