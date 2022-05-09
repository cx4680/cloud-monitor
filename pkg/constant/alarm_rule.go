package constant

//告警规则类型：单指标，多指标
const (
	AlarmRuleTypeSingleMetric   = 1
	AlarmRuleTypeMultipleMetric = 2
)

//告警规则逻辑组合：1:或, 2:且
const (
	AlarmRuleCombinationOr  = 1
	AlarmRuleCombinationAnd = 2
)
