package forms

// AlertRecordPageQueryForm 告警历史记录分页查询参数
type AlertRecordPageQueryForm struct {
	PageNum      int    `json:"pageNum"`
	PageSize     int    `json:"pageSize"`
	Region       string `json:"region"`
	Level        string `json:"level"`
	ResourceId   string `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	RuleId       string `json:"ruleId"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	RuleName     string `json:"ruleName"`
	Status       string `json:"status"`
}
