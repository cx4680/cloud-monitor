package model

import "time"

type AlarmInstance struct {
	InstanceID   string    `gorm:"column:instance_id;primary_key"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
	RegionCode   string    `gorm:"column:region_code"`
	ZoneCode     string    `gorm:"column:zone_code"`
	Ip           string    `gorm:"column:ip"`
	RegionName   string    `gorm:"column:region_name"`
	ZoneName     string    `gorm:"column:zone_name"`
	InstanceName string    `gorm:"column:instance_name"`
	TenantID     string    `gorm:"column:tenant_id"` // 租户id
	ProductType  string    `gorm:"column:product_type"`
}

func (*AlarmInstance) TableName() string {
	return "t_alarm_instance"
}
