package dto

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
	ResourceId         string   `json:"resourceId"`
	ResourceGroupId    string   `json:"resourceGroupId"`
}

type UserContactInfo struct {
	ContactId   string `gorm:"COLUMN:contactId"`
	ContactName string `gorm:"COLUMN:contactName"`
	Phone       string `gorm:"COLUMN:phone"`
	Mail        string `gorm:"COLUMN:mail"`
}

type ContactGroupInfo struct {
	GroupId   string            `gorm:"COLUMN:groupId"`
	GroupName string            `gorm:"COLUMN:groupName"`
	Contacts  []UserContactInfo `gorm:"-"`
}

type AutoScalingData struct {
	TenantId        string
	RuleId          string
	ResourceGroupId string
	Param           string
}
