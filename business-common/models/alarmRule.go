package models

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"time"
)

type AlarmRule struct {
	ID             string               `gorm:"column:id;primary_key"`
	MonitorType    string               `gorm:"column:monitor_type"` // 监控类型
	ProductType    string               `gorm:"column:product_type"` // 所属产品
	Dimensions     int                  `gorm:"column:dimensions"`   // 维度（1 全部资源 2 实例 ）
	Name           string               `gorm:"column:name"`         // 规则名称
	MetricName     string               `gorm:"column:metric_name"`
	RuleCondition  *forms.RuleCondition `gorm:"column:trigger_condition"`                             // 触发条件
	SilencesTime   string               `gorm:"column:silences_time"`                                 // 冷却周期
	EffectiveStart string               `gorm:"column:effective_start"`                               // 监控时间段-开始时间
	EffectiveEnd   string               `gorm:"column:effective_end"`                                 // 监控时间段-结束时间
	Level          int                  `gorm:"column:level"`                                         // 报警级别  紧急1 重要2次要3提醒 4
	NotifyChannel  int                  `gorm:"column:notify_channel"`                                // 通知方式 1 all  2 email  3 sms
	Enabled        int                  `gorm:"column:enabled;default:1"`                             // 启用（1）禁用（0）
	TenantID       string               `gorm:"column:tenant_id"`                                     // 租户id
	CreateTime     time.Time            `gorm:"column:create_time;autoCreateTime;default:time.now()"` // 创建时间
	CreateUser     string               `gorm:"column:create_user"`                                   // 创建人
	Deleted        int                  `gorm:"column:deleted;default:0"`                             // 未删除0  删除1
	UserName       string               `gorm:"column:user_name"`                                     // 租户名称
	UpdateTime     time.Time            `gorm:"column:update_time;autoCreateTime;default:time.now()"` // 创建时间
}

func (m *AlarmRule) TableName() string {
	return "t_alarm_rule"
}
