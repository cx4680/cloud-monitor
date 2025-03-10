package form

import (
	"encoding/json"
	"errors"
)

type AlarmPageReqParam struct {
	RuleName string `json:"ruleName,omitempty"`
	Status   string `json:"status,omitempty"` //AlarmRuleStatusEnum  enabled(1),disabled(0);
	TenantId string `json:"tenantId,omitempty"`
	PageSize int    `json:"pageSize" binding:"min=1,max=500"`
	Current  int    `json:"current" binding:"min=1"`
}

type AlarmRulePageDTO struct {
	Name           string       `json:"name" `
	MonitorType    string       `json:"monitorType" gorm:"monitor_type" `
	ProductType    string       `json:"productType" gorm:"productName" `
	InstanceNum    int          `json:"instanceNum" gorm:"column:instanceNum"`
	Status         string       `json:"status" gorm:"column:status"`
	RuleId         string       `json:"ruleId" gorm:"column:ruleId"`
	RuleConditions []*Condition `json:"ruleConditions" gorm:"-"`
}

type AlarmRuleDetailDTO struct {
	MonitorType      string          `json:"monitorType"`
	ProductType      string          `json:"productType"`
	Scope            string          `json:"scope"`
	InstanceList     []*InstanceInfo `json:"instanceList" gorm:"-"`
	RuleName         string          `json:"ruleName" gorm:"column:ruleName"`
	SilencesTime     string          `json:"silencesTime"`
	EffectiveStart   string          `json:"effectiveStart"  gorm:"column:effectiveStart"`
	EffectiveEnd     string          `json:"effectiveEnd"  gorm:"column:effectiveEnd"`
	AlarmLevel       int             `json:"alarmLevel" gorm:"column:alarmLevel"`
	TenantId         string          `json:"tenantId" gorm:"column:tenantId"`
	UserId           string          `json:"userId"`
	Id               string          `json:"id"`
	RuleConditions   []*Condition    `json:"ruleConditions" gorm:"-"`
	Status           string          `json:"status"`
	NoticeGroups     []*NoticeGroup  `json:"noticeGroups" gorm:"-"`
	Describe         string          `json:"describe" gorm:"column:describe"`
	AlarmHandlerList []*Handler      `json:"alarmHandlerList" gorm:"-"`
	Type             uint8           `json:"type"`
	Combination      uint8           `json:"combination"`
	Period           int             `json:"period"`
	Times            int             `json:"times"`
}
type NoticeGroup struct {
	Id       string      `json:"id" gorm:"column:id"`
	Name     string      `json:"name" gorm:"column:name"`
	UserList []*UserInfo `json:"userList" gorm:"-"`
}
type UserInfo struct {
	Id       string `json:"id" gorm:"column:id"`
	Phone    string `json:"phone" gorm:"column:phone"`
	Email    string `json:"email" gorm:"column:email"`
	UserName string `json:"userName" gorm:"column:userName"`
}

type ResGroupInfo struct {
	CalcMode     int             `json:"calcMode"`
	ResGroupId   string          `json:"resGroupId"`
	ResGroupName string          `json:"resGroupName"`
	ResourceList []*InstanceInfo `json:"resourceList"`
}

type AlarmRuleAddReqDTO struct {
	MonitorType       string          `json:"monitorType"  binding:"required"`
	ProductType       string          `json:"productType"  binding:"required"`
	ProductId         int             `json:"productId" `
	Scope             string          `json:"scope"`
	TenantId          string          `json:"tenantId"`
	UserId            string          `json:"userId"`
	MetricCode        string          `json:"metricCode"`
	ResourceGroupList []*ResGroupInfo `gorm:"-"`
	ResourceList      []*InstanceInfo `json:"instanceList" gorm:"-"`
	AlarmHandlerList  []*Handler      `json:"alarmHandlerList" gorm:"-"`
	RuleName          string          `json:"ruleName" binding:"required"`
	SilencesTime      string          `json:"silencesTime"`
	Level             uint8           `json:"level"`
	GroupList         []string        `json:"groupList" gorm:"-"`
	Source            string          `json:"source"`
	SourceType        uint8           `json:"sourceType"`
	Id                string          `json:"id"`
	Type              uint8           `json:"type" binding:"oneof=1 2"`
	Combination       uint8           `json:"combination"`
	Period            int             `json:"period"`
	Times             int             `json:"times"`
	Conditions        []Condition     `json:"conditions" binding:"required" gorm:"-"`
	EffectiveStart    string          `json:"effectiveStart"`
	EffectiveEnd      string          `json:"effectiveEnd"`
	TemplateBizId     string
}

type Condition struct {
	MetricName         string  `json:"metricName"`
	MetricCode         string  `json:"metricCode"`
	Period             int     `json:"period"`
	Times              int     `json:"times"`
	Statistics         string  `json:"statistics"  binding:"required"`
	ComparisonOperator string  `json:"comparisonOperator"  binding:"required"`
	Threshold          float64 `json:"threshold" binding:"required"`
	Unit               string  `json:"unit"`
	Labels             string  `json:"labels"`
	Level              uint8   `json:"level"`
	SilencesTime       string  `json:"silencesTime"`
	Express            string  `json:"express"`
}

func (c *Condition) Scan(v interface{}) error {
	var err error
	switch vt := v.(type) {
	case string:
		err = json.Unmarshal([]byte(vt), &c)
	case []byte:
		err = json.Unmarshal(vt, &c)
	default:
		return errors.New("rule condition 转换错误")
	}
	return err
}

type RuleReqDTO struct {
	Id       string `json:"id"  binding:"required"`
	Status   string `json:"status"`
	TenantId string `json:"tenantId"`
}
type Handler struct {
	HandleType   int    `json:"handleType" gorm:"column:handle_type"`     // 1邮件；2 短信；3弹性
	HandleParams string `json:"handleParams" gorm:"column:handle_params"` //回调地址
}
