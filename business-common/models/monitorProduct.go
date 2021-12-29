package models

type MonitorProduct struct {
	Id           uint64 `gorm:"column:id;primary_key;autoIncrement" json:"id"` // ID
	Name         string `gorm:"column:name" json:"name"`                       // 监控产品名称
	Status       int    `gorm:"column:status" json:"status"`                   // 状态 0:停用 1:启用
	Description  string `gorm:"column:description" json:"description"`         // 描述
	CreateUser   string `gorm:"column:create_user" json:"createUser"`          // 创建人
	CreateTime   string `gorm:"column:create_time" json:"createTime"`          // 创建时间
	Route        string `gorm:"column:route" json:"route"`                     // 路由
	Cron         string `gorm:"column:cron" json:"cron"`
	Host         string `gorm:"column:host;size:500" json:"host"`
	PageUrl      string `gorm:"column:page_url;size:500" json:"pageUrl"`
	Abbreviation string `gorm:"column:abbreviation" json:"abbreviation"` //简称
}

func (*MonitorProduct) TableName() string {
	return "monitor_product"
}
