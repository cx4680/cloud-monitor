package dto

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
)

// Deprecated
type RuleExpress struct {
	RuleId         string `gorm:"column:ruleId"`
	RuleName       string `gorm:"column:ruleName"`
	ProductType    string
	MonitorType    string
	Level          uint8 `gorm:"column:level"`
	NoticeChannel  uint8 `gorm:"column:noticeChannel"`
	RuleConditions []form.Condition
	NoticeGroupIds []*form.NoticeGroup `gorm:"-"`
	ResGroupId     string              `gorm:"column:resource_group_id"`
	CalcMode       int                 `gorm:"column:calc_mode"`
	SilencesTime   string              `gorm:"column:silences_time"`
	ResourceId     string              `gorm:"column:resource_id"`
	TenantId       string
	SourceType     uint8
	Source         string
	Type           uint8
	Combination    uint8
	Period         int
	Times          int
}
