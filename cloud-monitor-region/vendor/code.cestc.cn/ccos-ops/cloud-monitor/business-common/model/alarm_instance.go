package model

import "time"

type AlarmInstance struct {
	Id           uint64    `gorm:"column:id;primary_key;autoIncrement"`
	InstanceID   string    `gorm:"column:instance_id;"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
	RegionCode   string    `gorm:"column:region_code"`
	ZoneCode     string    `gorm:"column:zone_code"`
	Ip           string    `gorm:"column:ip"`
	RegionName   string    `gorm:"column:region_name"`
	ZoneName     string    `gorm:"column:zone_name"`
	InstanceName string    `gorm:"column:instance_name"`
	TenantID     string    `gorm:"column:tenant_id"` // 租户id
	ProductName  string    `gorm:"column:product_name"`
	ProductBizId string    `gorm:"column:product_biz_id"`
}

func (*AlarmInstance) TableName() string {
	return "t_resource"
}
