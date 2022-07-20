package form

// AlarmRecordPageQueryForm 告警历史记录分页查询参数
type AlarmRecordPageQueryForm struct {
	PageNum      int      `json:"pageNum" form:"pageNum"`
	PageSize     int      `json:"pageSize" form:"pageSize"`
	Region       string   `json:"region" form:"region"`
	Level        string   `json:"level" form:"level"`
	ResourceId   string   `json:"resourceId" form:"resourceId"`
	ResourceType []string `json:"resourceType" form:"resourceType"`
	RuleId       string   `json:"ruleId" form:"ruleId"`
	StartTime    string   `json:"startTime" form:"startTime"`
	EndTime      string   `json:"endTime" form:"endTime"`
	RuleName     string   `json:"ruleName" form:"ruleName"`
	Status       string   `json:"status" form:"status"`
	Expression   string   `json:"expression" form:"expression"`
	RegionCode   string   `json:"regionCode" form:"regionCode" `
	TenantID     string   `json:"tenantId" form:"tenantId"`
}
