package forms

import "container/list"

type PrometheusRequest struct {
	/**
	 * 租户Id（位于request head中）
	 */
	TenantId      string `form:"tenantId"`
	MonitorItemId string `form:"monitorItemId"`
	Name          string `form:"name"`
	Instance      string `form:"instance"`
	Labels        string `form:"labels"`
	/**
	 * 范围查询 s(秒)  m(分钟)  h(小时)  d(天)  w(周)  y(年)
	 */
	Scope string `form:"range"`
	/**
	 * 瞬时数据查询参数 时间戳
	 */
	Time string `form:"time"`
	/**
	 * 区间数据查询参数 时间戳
	 */
	Start int64 `form:"start"`
	End   int64 `form:"end"`
	Step  int64 `form:"step"`
	/**
	 * 统计方式
	 * 聚合函数 sum(求和)  min(最小值)  max (最大值)  avg (平均值)  stddev (标准差)  stdvar (标准差异)  count (计数)
	 */
	Statistics string `form:"statistics"`
}

type PrometheusResponse struct {
	Status string
	Data   PrometheusData
}

type PrometheusData struct {
	ResultType string
	Result     []PrometheusResult
}

type PrometheusResult struct {
	Metric map[string]string
	Value  []interface{}
	Values [][]interface{}
}

type PrometheusValue struct {
	Time  string
	Value string
}

type PrometheusAxis struct {
	XAxis list.List
	YAxis map[string][]string
}
