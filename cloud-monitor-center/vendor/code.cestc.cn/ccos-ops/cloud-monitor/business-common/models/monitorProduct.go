package models

type MonitorProduct struct {
	Id          string `gorm:"column:id;primary_key" json:"id"`       // ID
	Name        string `gorm:"column:name" json:"name"`               // 监控产品名称
	Status      int    `gorm:"column:status" json:"status"`           // 状态 0:停用 1:启用
	Description string `gorm:"column:description" json:"description"` // 描述
	CreateUser  string `gorm:"column:create_user" json:"create_user"` // 创建人
	CreateTime  string `gorm:"column:create_time" json:"create_time"` // 创建时间
	Route       string `gorm:"column:route" json:"route"`             // 路由
	Cron        string `gorm:"column:cron" json:"cron"`
	Host        string `gorm:"column:host;size:500" json:"host"`
	PageUrl     string `gorm:"column:pageUrl;size:500" json:"pageUrl"`
}

func (*MonitorProduct) TableName() string {
	return "monitor_product"
}
