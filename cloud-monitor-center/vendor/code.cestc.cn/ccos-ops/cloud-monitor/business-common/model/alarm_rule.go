package model

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"time"
)

type AlarmRule struct {
	Id             uint64              `gorm:"column:id;primary_key;autoIncrement"`
	BizId          string              `gorm:"column:biz_id;"`
	MonitorType    string              `gorm:"column:monitor_type"`   // 监控类型
	ProductName    string              `gorm:"column:product_name"`   // 所属产品
	ProductBizId   string              `gorm:"column:product_biz_id"` // 所属产品
	Dimensions     uint8               `gorm:"column:dimensions"`     // 维度（1 全部资源 2 实例 ）
	Name           string              `gorm:"column:name"`           // 规则名称
	MetricName     string              `gorm:"column:metric_name"`
	RuleCondition  *form.RuleCondition `gorm:"column:trigger_condition"`                        // 触发条件
	SilencesTime   string              `gorm:"column:silences_time"`                            // 冷却周期
	EffectiveStart string              `gorm:"column:effective_start"`                          // 监控时间段-开始时间
	EffectiveEnd   string              `gorm:"column:effective_end"`                            // 监控时间段-结束时间
	Level          uint8               `gorm:"column:level;"`                                   // 报警级别  紧急1 重要2次要3提醒 4
	Enabled        uint8               `gorm:"column:enabled;default:1;force"`                  // 启用（1）禁用（0）
	TenantID       string              `gorm:"column:tenant_id"`                                // 租户id
	CreateTime     time.Time           `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
	CreateUser     string              `gorm:"column:create_user"`                              // 创建人
	Deleted        uint8               `gorm:"column:deleted;default:0"`                        // 未删除0  删除1
	UserName       string              `gorm:"column:user_name"`                                // 租户名称
	UpdateTime     time.Time           `gorm:"column:update_time;autoUpdateTime;type:datetime"` // 创建时间
	Source         string              `gorm:"column:source;"`                                  // 请求来源（如弹性伸缩组id）
	SourceType     uint8               `gorm:"column:source_type;default:1;force"`              //1 页面 ,2弹性伸缩
}

func (m *AlarmRule) TableName() string {
	return "t_alarm_rule"
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
