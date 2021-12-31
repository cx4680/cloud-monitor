package vo

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/util"
	"fmt"
	"strings"
)

type InstanceRuleVO struct {
	Id          string
	Name        string
	ProductType string
	MonitorType string
	MonitorItem string
	Condition   string
}

type RuleCondition struct {
	MetricName         string
	Unit               string
	Labels             string
	MonitorItemName    string
	Statistics         string
	ComparisonOperator string
	Threshold          string
	Period             string
	Times              int
}

func (ruleCondition *RuleCondition) GetExpress() string {
	uni := util.If(ruleCondition.Unit == "" || strings.EqualFold(ruleCondition.Unit, "null"), "", ruleCondition.Unit)
	return fmt.Sprintf("%s%s%s%s%s 统计周期%s 持续%s个周期", ruleCondition.MonitorItemName, ruleCondition.Statistics, ruleCondition.ComparisonOperator, ruleCondition.Threshold, uni, ruleCondition.Period, ruleCondition.Times)
}
