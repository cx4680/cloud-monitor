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
	MonitorItem   string         `json:"monitorItem" gorm:"monitor_item" `
	Express       string         `json:"express" gorm:"express"`
	InstanceNum   int            `json:"instanceNum" gorm:"column:instanceNum"`
	Status        string         `json:"status" gorm:"column:status"`
	RuleId        string         `json:"ruleId" gorm:"column:ruleId"`
	RuleCondition *RuleCondition `json:"ruleCondition" gorm:"column:trigger_condition"`
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
	RuleCondition    *RuleCondition  `json:"ruleCondition" gorm:"column:ruleCondition"`
	Status           string          `json:"status"`
	NoticeGroups     []*NoticeGroup  `json:"noticeGroups" gorm:"-"`
	Describe         string          `json:"describe" gorm:"column:describe"`
	AlarmHandlerList []*Handler      `json:"alarmHandlerList" gorm:"-"`
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
	MonitorType       string `json:"monitorType"  binding:"required"`
	ProductType       string `json:"productType"  binding:"required"`
	ProductId         int    `json:"productId" `
	Scope             string `json:"scope"`
	TenantId          string `json:"tenantId"`
	UserId            string `json:"userId"`
	ResourceGroupList []*ResGroupInfo
	ResourceList      []*InstanceInfo `json:"instanceList"`
	AlarmHandlerList  []*Handler      `json:"alarmHandlerList"`
	RuleName          string          `json:"ruleName"`
	RuleCondition     *RuleCondition  `json:"triggerCondition"`
	SilencesTime      string          `json:"silencesTime"`
	AlarmLevel        int             `json:"alarmLevel"`
	NoticeChannel     string          `json:"noticeChannel"`
	GroupList         []string        `json:"groupList"`
	Source            string          `json:"source"`
	SourceType        int             `json:"sourceType"`
	Id                string          `json:"id"`
}

type RuleCondition struct {
	MetricName         string  `json:"metricName"`
	Period             int     `json:"period"`
	Times              int     `json:"times"`
	Statistics         string  `json:"statistics"`
	ComparisonOperator string  `json:"comparisonOperator"`
	Threshold          float64 `json:"threshold"`
	Unit               string  `json:"unit"`
	Labels             string  `json:"labels"`
	MonitorItemName    string  `json:"monitorItemName"`
}
type RuleReqDTO struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	TenantId string `json:"tenantId"`
}
type Handler struct {
	HandleType   int    `json:"handleType" gorm:"column:handle_type"`   // 1邮件；2 短信；3弹性
	HandleParams string `json:"handleParams" gorm:"column:handle_params"` //回调地址
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
