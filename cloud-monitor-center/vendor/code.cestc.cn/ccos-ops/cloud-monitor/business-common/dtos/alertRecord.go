package dtos

type RuleDesc struct {
	RuleName           string   `json:"rule_name"`
	Product            string   `json:"product"`
	InstanceInfo       string   `json:"instance_info"`
	MetricName         string   `json:"metric_name"`
	ComparisonOperator string   `json:"comparison_operator"`
	TargetValue        int      `json:"target_value"`
	Time               int      `json:"time"`
	Period             int      `json:"period"`
	Unit               string   `json:"unit"`
	Express            string   `json:"express"`
	Level              int      `json:"level"`
	NotifyChannel      string   `json:"notify_channel"`
	MonitorItem        string   `json:"monitor_item"`
	MonitorType        string   `json:"monitor_type"`
	InstanceId         string   `json:"instance_id"`
	RuleId             string   `json:"rule_id"`
	TenantId           string   `json:"tenant_id"`
	RegionName         string   `json:"region_name"`
	Statistic          string   `json:"statistic"`
	GroupList          []string `json:"group_list"`
	ResourceGroupId    string   `json:"ResourceGroupId"`
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
