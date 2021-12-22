package dtos

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
)

type RuleExpress struct {
	RuleId         string `gorm:"column:ruleId"`
	RuleName       string `gorm:"column:ruleName"`
	ProductType    string
	MonitorType    string
	Level          int                  `gorm:"column:level"`
	NoticeChannel  int                  `gorm:"column:noticeChannel"`
	RuleCondition  *forms.RuleCondition `gorm:"column:ruleCondition"`
	NoticeGroupIds []*forms.NoticeGroup `gorm:"-"`
	ResGroupId     string               `gorm:"column:resource_group_id"`
	CalcMode       int                  `gorm:"column:calc_mode"`
	SilencesTime   string               `gorm:"column:silences_time"`
	ResourceId     string               `gorm:"column:resource_id"`
	TenantId       string
	SourceType     int
}
