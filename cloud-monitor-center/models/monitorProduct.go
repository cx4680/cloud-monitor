package models

import "time"

type MonitorProduct struct {
	ID          string    `gorm:"column:id;primary_key" json:"id"`       // ID
	Name        string    `gorm:"column:name" json:"name"`               // 监控产品名称
	Status      int       `gorm:"column:status" json:"status"`           // 状态 0:停用 1:启用
	Description string    `gorm:"column:description" json:"description"` // 描述
	CreateUser  string    `gorm:"column:create_user" json:"create_user"` // 创建人
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"` // 创建时间
}

func (*MonitorProduct) TableName() string {
	return "monitor_product"
}
