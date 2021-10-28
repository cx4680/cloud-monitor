package models

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"time"
)

type AlarmRule struct {
	Id             string               `orm:"id" json:"id"`
	MonitorType    string               `orm:"monitor_type" json:"monitor_type"` // 监控类型
	ProductType    string               `orm:"product_type" json:"product_type"` // 所属产品
	Dimensions     int                  `orm:"dimensions" json:"dimensions"`     // 维度（1 全部资源 2 实例 ）
	Name           string               `orm:"name" json:"name"`                 // 规则名称
	MetricName     string               `orm:"metric_name" json:"metric_name"`
	RuleCondition  *forms.RuleCondition `orm:"trigger_condition" json:"trigger_condition"`                               // 触发条件
	SilencesTime   string               `orm:"silences_time" json:"silences_time"`                                       // 冷却周期
	EffectiveStart string               `orm:"effective_start" json:"effective_start"`                                   // 监控时间段-开始时间
	EffectiveEnd   string               `orm:"effective_end" json:"effective_end"`                                       // 监控时间段-结束时间
	Level          int                  `orm:"level" json:"level"`                                                       // 报警级别
	NotifyChannel  int                  `orm:"notify_channel" json:"notify_channel"`                                     // 通知方式 1 all  2 email  3 sms
	Enabled        int                  `orm:"enabled" json:"enabled"`                                                   // 启用（1）禁用（0）
	TenantId       string               `orm:"tenant_id" json:"tenant_id"`                                               // 租户id
	CreateTime     time.Time            `gorm:"column:create_time;autoCreateTime;default:time.now()" json:"create_time"` // 创建时间
	CreateUser     string               `orm:"create_user" json:"create_user"`                                           // 创建人
	Deleted        int                  `orm:"deleted" json:"deleted"`                                                   // 未删除0  删除1
	UserName       string               `orm:"user_name" json:"user_name"`                                               // 租户名称
	UpdateTime     time.Time            `gorm:"column:update_time;autoCreateTime;default:time.now()" json:"update_time"` // 创建时间
}

func (*AlarmRule) TableName() string {
	return "t_alarm_rule"
}
