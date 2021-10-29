package models

import "time"

type MonitorItem struct {
	ID             int       `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // ID
	ProductID      string    `gorm:"column:product_id"`                                    // 监控产品ID
	Name           string    `gorm:"column:name"`                                          // 监控项名称
	MetricName     string    `gorm:"column:metric_name"`                                   // 指标名
	Labels         string    `gorm:"column:labels"`                                        // 标签名 分号隔开
	MetricsLinux   string    `gorm:"column:metrics_linux"`                                 // linux表达式
	MetricsWindows string    `gorm:"column:metrics_windows"`                               // windows表达式
	Statistics     string    `gorm:"column:statistics"`                                    // 统计方式
	Unit           string    `gorm:"column:unit"`                                          // 单位
	Frequency      string    `gorm:"column:frequency"`                                     // 频率
	Type           int       `gorm:"column:type"`                                          // 类型 1:基础监控 2:操作系统监控
	IsDisplay      int       `gorm:"column:is_display"`                                    // 是否显示 0:不显示 1:显示
	Status         int       `gorm:"column:status"`                                        // 状态 0:停用 1:启用
	Description    string    `gorm:"column:description"`                                   // 描述
	CreateUser     string    `gorm:"column:create_user"`                                   // 创建人
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime;default:time.now()"` // 创建时间
}

func (*MonitorItem) TableName() string {
	return "monitor_item"
}
