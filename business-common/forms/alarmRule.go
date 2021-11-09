package forms

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

type AlarmPageReqParam struct {
	RuleName string `json:"ruleName,omitempty"`
	Status   string `json:"status,omitempty"` //AlarmRuleStatusEnum  enabled(1),disabled(0);
	TenantId string `json:"tenantId,omitempty"`
	PageSize int    `json:"pageSize" binding:"min=1,max=500"`
	Current  int    `json:"current" binding:"min=1"`
}

type AlarmRulePageDTO struct {
	Name          string         `json:"name" `
	MonitorType   string         `json:"monitorType" gorm:"monitor_type" `
	ProductType   string         `json:"productType" gorm:"product_type" `
	MetricName    string         `json:"metricName" gorm:"metric_name" `
	Express       string         `json:"express" gorm:"express"`
	InstanceNum   int            `json:"instanceNum" gorm:"column:instanceNum"`
	Status        string         `json:"status" gorm:"column:status"`
	RuleId        string         `json:"ruleId" gorm:"column:ruleId"`
	RuleCondition *RuleCondition `json:"ruleCondition" gorm:"column:trigger_condition"`
}

type AlarmRuleDetailDTO struct {
	MonitorType       string          `json:"monitorType"`
	ProductType       string          `json:"productType"`
	Scope             string          `json:"scope"`
	InstanceList      []*InstanceInfo `json:"instanceList" gorm:"-"`
	RuleName          string          `json:"ruleName" gorm:"column:ruleName"`
	SilencesTime      string          `json:"silencesTime"`
	EffectiveStart    string          `json:"effectiveStart"  gorm:"column:effectiveStart"`
	EffectiveEnd      string          `json:"effectiveEnd"  gorm:"column:effectiveEnd"`
	AlarmLevel        int             `json:"alarmLevel" gorm:"column:alarmLevel"`
	NoticeChannel     string          `json:"noticeChannel" gorm:"column:noticeChannel"`
	TenantId          string          `json:"tenantId" gorm:"column:tenantId"`
	UserId            string          `json:"userId"`
	Id                string          `json:"id"`
	RuleCondition     *RuleCondition  `json:"ruleCondition" gorm:"column:ruleCondition"`
	Status            string          `json:"status"`
	NoticeGroups      []*NoticeGroup  `json:"noticeGroups" gorm:"-"`
	NoticeChannelDesc string          `json:"noticeChannelDesc" gorm:"column:noticeChannelDesc"`
	Describe          string          `json:"describe" gorm:"column:describe"`
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

type AlarmRuleAddReqDTO struct {
	MonitorType   string          `json:"monitorType"`
	ProductType   string          `json:"productType"`
	Scope         string          `json:"scope"`
	TenantId      string          `json:"tenantId"`
	UserId        string          `json:"userId"`
	InstanceList  []*InstanceInfo `json:"instanceList"`
	RuleName      string          `json:"ruleName"`
	RuleCondition *RuleCondition  `json:"triggerCondition"`
	SilencesTime  string          `json:"silencesTime"`
	AlarmLevel    int             `json:"alarmLevel"`
	NoticeChannel string          `json:"noticeChannel"`
	GroupList     []string        `json:"groupList"`
	Id            string          `json:"id"`
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
	Id       string `json:"id"`
	Status   string `json:"status"`
	TenantId string `json:"tenantId"`
}

func (p *RuleCondition) Value() (driver.Value, error) {
	bs, err := json.Marshal(p)
	return string(bs), errors.WithStack(err)
}
func (s *RuleCondition) Scan(v interface{}) error {
	var err error
	switch vt := v.(type) {
	case string:
		err = json.Unmarshal([]byte(vt), &s)
	case []byte:
		err = json.Unmarshal(vt, &s)
	default:
		return errors.New("rule condition 转换错误")
	}
	return err
}
