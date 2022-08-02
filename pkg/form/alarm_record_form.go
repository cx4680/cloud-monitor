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
	TenantId     string   `json:"tenantId" form:"tenantId"`
	IamUserId    string   `json:"iamUserId" form:"iamUserId"`
}

type AlarmRecordNum struct {
	P1 int `json:"p1" gorm:"column:p1"`
	P2 int `json:"p2" gorm:"column:p2"`
	P3 int `json:"p3" gorm:"column:p3"`
	P4 int `json:"p4" gorm:"column:p4"`
}

type IamDirectory struct {
	Module struct {
		DirectoryId int `json:"directoryId"`
		ChildList   []*struct {
			DirectoryId int `json:"directoryId"`
		} `json:"childList"`
	} `json:"module"`
}
