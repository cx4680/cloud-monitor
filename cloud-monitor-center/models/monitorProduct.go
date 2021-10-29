package models

import "time"

type MonitorProduct struct {
	ID          string    `gorm:"column:id;primary_key"` // ID
	Name        string    `gorm:"column:name"`           // 监控产品名称
	Status      int       `gorm:"column:status"`         // 状态 0:停用 1:启用
	Description string    `gorm:"column:description"`    // 描述
	CreateUser  string    `gorm:"column:create_user"`    // 创建人
	CreateTime  time.Time `gorm:"column:create_time"`    // 创建时间
}

func (*MonitorProduct) TableName() string {
	return "monitor_product"
}
