package form

import "time"

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
	RegionCode   string   `json:"regionCode" form:"regionCode"`
	ProductCode  string   `json:"productCode" form:"productCode"`
	TenantId     string   `json:"tenantId" form:"tenantId"`
	IamUserId    string   `json:"iamUserId" form:"iamUserId"`
}

type AlarmRecordNum struct {
	Level int `json:"level" gorm:"column:level"`
	Count int `json:"count" gorm:"column:count"`
}

type ProductAlarmRecordNum struct {
	ProductCode string `json:"productCode" gorm:"column:productCode"`
	Count       int    `json:"count" gorm:"column:count"`
}

type AlarmRecordPage struct {
	ProductCode string    `json:"productCode" gorm:"column:productCode"`
	InstanceId  string    `json:"instanceId" gorm:"column:instanceId"`
	RuleName    string    `json:"ruleName" gorm:"column:ruleName"`
	Level       string    `json:"level" gorm:"column:level"`
	Time        time.Time `json:"-" gorm:"column:time"`
	FmtTime     string    `json:"time"`
}
