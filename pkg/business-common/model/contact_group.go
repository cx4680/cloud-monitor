package model

import "time"

type ContactGroup struct {
	Id          uint64    `gorm:"column:id;primary_key;autoIncrement" json:"id"` // ID
	BizId       string    `gorm:"column:biz_id" json:"bizId"`                    // 关联ID
	TenantId    string    `gorm:"column:tenant_id" json:"tenant_id"`             // 租户ID
	Name        string    `gorm:"column:name" json:"name"`                       // 组名
	Description string    `gorm:"column:description" json:"description"`         // 描述
	CreateUser  string    `gorm:"column:create_user" json:"create_user"`         // 创建人
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`         // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`         // 修改时间
	State       uint8     `gorm:"column:state" json:"state"`                     // 状态1启动2删除
}

func (*ContactGroup) TableName() string {
	return "t_contact_group"
}
