package dtos

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
)

type RuleExpress struct {
	RuleId         string `gorm:"column:ruleId"`
	RuleName       string `gorm:"column:ruleName"`
	ProductType    string
	MonitorType    string
	Level          int                  `gorm:"column:level"`
	NoticeChannel  int                  `gorm:"column:noticeChannel"`
	RuleCondition  *forms.RuleCondition `gorm:"column:ruleCondition"`
	NoticeGroupIds []*forms.NoticeGroup `gorm:"-"`
	ResGroupId     string               `gorm:"column:resource_group_id"`
	CalcMode       int                  `gorm:"column:calc_mode"`
	SilencesTime   string               `gorm:"column:silences_time"`
	ResourceId     string               `gorm:"column:resource_id"`
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
