package model

type MonitorProduct struct {
	Id           uint64 `gorm:"column:id;primary_key;autoIncrement" json:"id"`   // ID
	BizId        string `gorm:"column:biz_id;size=50" json:"bizId"`              // 关联ID
	Name         string `gorm:"column:name" json:"name"`                         // 监控产品名称
	Status       uint8  `gorm:"column:status" json:"status"`                     // 状态 0:停用 1:启用
	Description  string `gorm:"column:description" json:"description"`           // 描述
	CreateUser   string `gorm:"column:create_user" json:"createUser"`            // 创建人
	CreateTime   string `gorm:"column:create_time" json:"createTime"`            // 创建时间
	Route        string `gorm:"column:route" json:"route"`                       // 路由
	Cron         string `gorm:"column:cron" json:"cron"`                         // 定时任务
	Host         string `gorm:"column:host;size:500" json:"host"`                // 请求路径svc
	PageUrl      string `gorm:"column:page_url;size:500" json:"pageUrl"`         // 请求路径
	Abbreviation string `gorm:"column:abbreviation;size=20" json:"abbreviation"` // 简称
	Sort         uint64 `gorm:"column:sort" json:"sort"`                         // 排序
	MonitorType  string `gorm:"column:monitor_type" json:"monitorType"`          // 监控类型
}

func (*MonitorProduct) TableName() string {
	return "t_monitor_product"
}

type MonitorProductDTO struct {
	Id           uint64 `gorm:"column:id;primary_key;autoIncrement" json:"id"`   // ID
	BizId        string `gorm:"column:biz_id;size=50" json:"bizId"`              // 关联ID
	Name         string `gorm:"column:name" json:"name"`                         // 监控产品名称
	Status       uint8  `gorm:"column:status" json:"status"`                     // 状态 0:停用 1:启用
	Description  string `gorm:"column:description" json:"description"`           // 描述
	CreateUser   string `gorm:"column:create_user" json:"createUser"`            // 创建人
	CreateTime   string `gorm:"column:create_time" json:"createTime"`            // 创建时间
	Route        string `gorm:"column:route" json:"route"`                       // 路由
	Abbreviation string `gorm:"column:abbreviation;size=20" json:"abbreviation"` // 简称
	Sort         uint64 `gorm:"column:sort" json:"sort"`                         // 排序
	MonitorType  string `gorm:"column:monitor_type" json:"monitorType"`          // 监控类型
}
