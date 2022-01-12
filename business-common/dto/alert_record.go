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
	Level              uint8    `json:"level"`
	MonitorItem        string   `json:"monitorItem"`
	MonitorType        string   `json:"monitorType"`
	RuleId             string   `json:"ruleId"`
	TenantId           string   `json:"tenantId"`
	RegionName         string   `json:"regionName"`
	Statistic          string   `json:"statistic"`
	GroupList          []string `json:"groupList"`
	ResourceId         string   `json:"resourceId"`
	ResourceGroupId    string   `json:"resourceGroupId"`
	Source             string   `json:"source"`     // 请求来源（如弹性伸缩组id）
	SourceType         uint8    `json:"sourceType"` //1 页面 ,2弹性伸缩
}

type UserContactInfo struct {
	ContactId   string `gorm:"COLUMN:contactId" json:"contactId"`
	ContactName string `gorm:"COLUMN:contactName" json:"contactName"`
	Phone       string `gorm:"COLUMN:phone" json:"phone"`
	Mail        string `gorm:"COLUMN:mail" json:"mail"`
}

type ContactGroupInfo struct {
	GroupId   string            `gorm:"COLUMN:groupId" json:"groupId"`
	GroupName string            `gorm:"COLUMN:groupName" json:"groupName"`
	Contacts  []UserContactInfo `gorm:"-" json:"contacts"`
}

type AutoScalingData struct {
	TenantId        string
	RuleId          string
	ResourceGroupId string
	Param           string
}
