package models

type MonitorItem struct {
	Id             string `orm:"id" json:"id"`                          // ID
	ProductId      string `orm:"product_id" json:"productId"`           // 监控产品ID
	Name           string `orm:"name" json:"name"`                      // 监控项名称
	MetricName     string `orm:"metric_name" json:"metricName"`         // 指标名
	Labels         string `orm:"labels" json:"labels"`                  // 指标名
	MetricsWindows string `orm:"metrics_windows" json:"metricsWindows"` // windows表达式
	Statistics     string `orm:"statistics" json:"statistics"`          // 统计方式
	Unit           string `orm:"unit" json:"unit"`                      // 单位
	Frequency      string `orm:"frequency" json:"frequency"`            // 频率
	Type           int    `orm:"type" json:"type"`                      // 类型 1:基础监控 2:操作系统监控
	IsDisplay      int    `orm:"is_display" json:"isDisplay"`           // 是否显示 0:不显示 1:显示
	Status         int    `orm:"status" json:"status"`                  // 状态 0:停用 1:启用
	Description    string `orm:"description" json:"description"`        // 描述
	CreateUser     string `orm:"create_user" json:"createUser"`         // 创建时间
	CreateTime     string `orm:"create_time" json:"createTime"`         // 创建人
}

func (*MonitorItem) TableName() string {
	return "monitor_item"
}
