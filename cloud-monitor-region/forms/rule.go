package forms

type TriggerCondition struct {
	MetricName string
	/**
	 * 周期
	 */
	Period int
	/*
	 * 次数
	 */
	Times int
	/**
	 * 统计类型
	 */
	Statistics string
	/**
	 * > >= < <= ==
	 */
	ComparisonOperator string
	/**
	 * 门限值
	 */
	Threshold int
}
