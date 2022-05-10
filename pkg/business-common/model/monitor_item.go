package model

import "time"

type MonitorItem struct {
	Id             uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                    // ID
	BizId          string    `gorm:"column:biz_id;size=50" json:"bizId"`                                // 关联ID
	ProductBizID   string    `gorm:"column:product_biz_id" json:"productBizId"`                         // 监控产品ID
	Name           string    `gorm:"column:name" json:"name"`                                           // 监控项名称
	MetricName     string    `gorm:"column:metric_name" json:"metricName"`                              // 指标名
	Labels         string    `gorm:"column:labels" json:"labels"`                                       // 标签名 分号隔开
	MetricsLinux   string    `gorm:"column:metrics_linux" json:"metricsLinux"`                          // linux表达式
	MetricsWindows string    `gorm:"column:metrics_windows" json:"metricsWindows"`                      // windows表达式
	Statistics     string    `gorm:"column:statistics" json:"statistics"`                               // 统计方式
	Unit           string    `gorm:"column:unit" json:"unit"`                                           // 单位
	Frequency      string    `gorm:"column:frequency" json:"frequency"`                                 // 频率
	Type           uint8     `gorm:"column:type" json:"type"`                                           // 类型 1:基础监控 2:操作系统监控
	IsDisplay      uint8     `gorm:"column:is_display" json:"isDisplay"`                                // 是否显示 0:不显示 1:显示
	Display        string    `gorm:"column:display" json:"display"`                                     // 展示位置
	Status         uint8     `gorm:"column:status" json:"status"`                                       // 状态 0:停用 1:启用
	Description    string    `gorm:"column:description" json:"description"`                             // 描述
	CreateUser     string    `gorm:"column:create_user" json:"createUser"`                              // 创建人
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime;type:datetime" json:"createTime"` // 创建时间
	ShowExpression string    `gorm:"column:show_expression;autoCreateTime" json:"showExpression"`       // 展示表达式
}

func (*MonitorItem) TableName() string {
	return "t_monitor_item"
}
