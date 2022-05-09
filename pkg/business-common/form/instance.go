package form

type InstanceRulePageReqParam struct {
	InstanceId string `json:"instanceId" binding:"required"`
	PageSize   int    `json:"pageSize" binding:"min=1,max=500"`
	Current    int    `json:"current" binding:"min=1"`
}

type InstanceRuleDTO struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	ProductType string   `json:"productType"`
	MonitorType string   `json:"monitorType"`
	Condition   []string `json:"condition"`
}

type InstanceInfo struct {
	InstanceId   string `json:"instanceId"  binding:"required"`
	ZoneCode     string `json:"zoneCode"`
	RegionCode   string `json:"regionCode" binding:"required"`
	RegionName   string `json:"regionName" binding:"required"`
	ZoneName     string `json:"zoneName"`
	Ip           string `json:"ip"`
	Status       string `json:"status"`
	InstanceName string `json:"instanceName" binding:"required"`
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
	BindRuleList   []*InstanceRuleDTO `json:"bindRuleList"`
	UnbindRuleList []*InstanceRuleDTO `json:"unbindRuleList"`
	TenantId       string
}

type UnBindRuleParam struct {
	InstanceId string `json:"instanceId" binding:"required"`
	RuleId     string `json:"ruleId" binding:"required"`
	TenantId   string
}
