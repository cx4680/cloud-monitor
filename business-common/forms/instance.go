package forms

type InstanceRulePageReqParam struct {
	InstanceId string `json:"instanceId" binding:"required"`
	PageSize   int    `json:"pageSize" binding:"min=1,max=500"`
	Current    int    `json:"current" binding:"min=1"`
}

type InstanceRuleDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProductType string `json:"productType"`
	MonitorType string `json:"monitorType"`
	MonitorItem string `json:"monitorItem" gorm:"column:monitorItem"`
	Condition   string `json:"condition" gorm:"column:ruleCondition"`
}

type InstanceInfo struct {
	InstanceId   string `json:"instanceId"`
	ZoneCode     string `json:"zoneCode"`
	RegionCode   string `json:"regionCode"`
	RegionName   string `json:"regionName"`
	ZoneName     string `json:"zoneName"`
	Ip           string `json:"ip"`
	Status       string `json:"status"`
	InstanceName string `json:"instanceName"`
}

type InstanceBindRuleDTO struct {
	TenantId string
	InstanceInfo
	RuleIdList []string `json:"ruleIdList"`
}

type ProductRuleParam struct {
	MonitorType string `json:"monitorType" binding:"required"`
	ProductType string `json:"productType" binding:"required"`
	InstanceId  string `json:"instanceId" binding:"required"`
	TenantId    string
}

type ProductRuleListDTO struct {
	BindRuleList   []InstanceRuleDTO `json:"bindRuleList"`
	UnbindRuleList []InstanceRuleDTO `json:"unbindRuleList"`
}

type UnBindRuleParam struct {
	InstanceId string `json:"instanceId" binding:"required"`
	RulId      string `json:"rulId" binding:"required"`
}
