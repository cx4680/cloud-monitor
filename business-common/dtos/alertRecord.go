package dtos

type RuleDesc struct {
	RuleName           string   `json:"ruleName"`
	Product            string   `json:"product"`
	MetricName         string   `json:"metricName"`
	ComparisonOperator string   `json:"comparisonOperator"`
	TargetValue        float64  `json:"targetValue"`
	Time               int      `json:"time"`
	Period             int      `json:"period"`
	Unit               string   `json:"unit"`
	Express            string   `json:"express"`
	Level              int      `json:"level"`
	MonitorItem        string   `json:"monitorItem"`
	MonitorType        string   `json:"monitorType"`
	RuleId             string   `json:"ruleId"`
	TenantId           string   `json:"tenantId"`
	RegionName         string   `json:"regionName"`
	Statistic          string   `json:"statistic"`
	GroupList          []string `json:"groupList"`
	ResourceGroupId    string   `json:"resourceGroupId"`
}

type UserContactInfo struct {
	ContactId   string
	ContactName string
	Phone       string
	Mail        string
}

type ContactGroupInfo struct {
	GroupId,
	GroupName string
	Contacts []UserContactInfo
}

type AutoScalingData struct {
	TenantId        string
	RuleId          string
	ResourceGroupId string
	Param           string
}
