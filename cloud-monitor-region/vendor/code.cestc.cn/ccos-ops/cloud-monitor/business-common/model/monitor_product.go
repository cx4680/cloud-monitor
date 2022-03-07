package model

type MonitorProduct struct {
	Id           uint64 `gorm:"column:id;primary_key;autoIncrement" json:"id"` // ID
	BizId        string `gorm:"column:biz_id;size=50" json:"bizId"`
	Name         string `gorm:"column:name" json:"name"`               // 监控产品名称
	Status       uint8  `gorm:"column:status" json:"status"`           // 状态 0:停用 1:启用
	Description  string `gorm:"column:description" json:"description"` // 描述
	CreateUser   string `gorm:"column:create_user" json:"createUser"`  // 创建人
	CreateTime   string `gorm:"column:create_time" json:"createTime"`  // 创建时间
	Route        string `gorm:"column:route" json:"route"`             // 路由
	Cron         string `gorm:"column:cron" json:"-"`
	Host         string `gorm:"column:host;size:500" json:"-"`
	PageUrl      string `gorm:"column:page_url;size:500" json:"-"`
	Abbreviation string `gorm:"column:abbreviation;size=20" json:"abbreviation"` //简称
	Sort         uint64 `gorm:"column:sort" json:"sort"`                         //排序
}

func (*MonitorProduct) TableName() string {
	return "t_monitor_product"
}
