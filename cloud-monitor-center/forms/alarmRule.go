package forms

type AlarmPageReqParam struct {
	RuleName string `json:"omitempty"`
	Status   string `json:"omitempty"` //AlarmRuleStatusEnum  enabled(1),disabled(0);
	TenantId string `json:",omitempty"`
	PageSize int    `json:",omitempty,default=10"`
	Current  int    `json:",omitempty,default=1"`
}

type AlarmRulePageDTO struct {
	Name        string `json:"name"`
	MonitorType string `json:"monitorType" orm:"monitor_type" `
	ProductType string `json:"productType" orm:"product_type" `
	MetricName  string `json:"metricName" orm:"metric_name" `
	Express     string `json:"express"`
	InstanceNum int    `json:"instanceNum"`
	Status      string `json:"status"`
	RuleId      string `json:"ruleId"`
}

type AlarmRuleDetailDTO struct {
	MonitorType      string         `json:"monitorType"`
	ProductType      string         `json:"productType"`
	Scope            string         `json:"scope"`
	InstanceList     []InstanceInfo `json:"instanceList"`
	RuleName         string         `json:"ruleName"`
	TriggerCondition interface{}    `json:"triggerCondition"`
	SilencesTime     string         `json:"silencesTime"`
	EffectiveStart   interface{}    `json:"effectiveStart"`
	EffectiveEnd     interface{}    `json:"effectiveEnd"`
	AlarmLevel       int            `json:"alarmLevel"`
	NoticeChannel    string         `json:"noticeChannel"`
	GroupList        interface{}    `json:"groupList"`
	TenantId         interface{}    `json:"tenantId"`
	UserId           interface{}    `json:"userId"`
	Id               string         `json:"id"`
	RuleCondition    RuleCondition  `json:"ruleCondition"`
	Status           string         `json:"status"`
	NoticeGroups     []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		UserList []struct {
			Phone    interface{} `json:"phone"`
			Email    string      `json:"email"`
			UserName string      `json:"userName"`
		} `json:"userList"`
	} `json:"noticeGroups"`
	NoticeChannelDesc string `json:"noticeChannelDesc"`
	Describe          string `json:"describe"`
}

type AlarmRuleAddReqDTO struct {
	MonitorType   string         `json:"monitorType"`
	ProductType   string         `json:"productType"`
	Scope         string         `json:"scope"`
	TenantId      string         `json:"tenantId"`
	UserId        string         `json:"userId"`
	InstanceList  []InstanceInfo `json:"instanceList"`
	RuleName      string         `json:"ruleName"`
	RuleCondition RuleCondition  `json:"triggerCondition"`
	SilencesTime  string         `json:"silencesTime"`
	AlarmLevel    int            `json:"alarmLevel"`
	NoticeChannel string         `json:"noticeChannel"`
	GroupList     []string       `json:"groupList"`
	Id            string         `json:"id"`
}

type RuleCondition struct {
	MetricName         string `json:"metricName"`
	Period             int    `json:"period"`
	Times              int    `json:"times"`
	Statistics         string `json:"statistics"`
	ComparisonOperator string `json:"comparisonOperator"`
	Threshold          int    `json:"threshold"`
	Unit               string `json:"unit"`
	Labels             string `json:"labels"`
	MonitorItemName    string `json:"monitorItemName"`
}
type RuleReqDTO struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
