package forms

type InstanceRulePageReqParam struct {
	InstanceId string `json:"instanceId"`
	PageSize   int    `json:"pageSize"`
	Current    int    `json:"current"`
}

type InstanceRuleDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProductType string `json:"productType"`
	MonitorType string `json:"monitorType"`
	MonitorItem string `json:"monitorItem"`
	Condition   string `json:"condition"`
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
	MonitorType string `json:"monitorType"`
	ProductType string `json:"productType"`
	InstanceId  string `json:"instanceId"`
}

type ProductRuleListDTO struct {
	BindRuleList   []BindRuleInfo `json:"bindRuleList"`
	UnbindRuleList []BindRuleInfo `json:"unbindRuleList"`
}

type BindRuleInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProductType string `json:"productType"`
	MonitorType string `json:"monitorType"`
	MonitorItem string `json:"monitorItem"`
	Condition   string `json:"condition"`
}
