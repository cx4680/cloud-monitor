package models

type MonitorProduct struct {
	Id          string `orm:"id" json:"id"`                   // ID
	Name        string `orm:"name" json:"name"`               // 监控产品名称
	Status      int    `orm:"status" json:"status"`           // 状态 0:停用 1:启用
	Description string `orm:"description" json:"description"` // 描述
	CreateUser  string `orm:"create_user" json:"create_user"` // 创建人
	//CreateTime  string `orm:"create_time" json:"create_time"` // 创建时间
}

func (*MonitorProduct) TableName() string {
	return "monitor_product"
}
