package model

import "time"

type Contact struct {
	Id          uint64    `gorm:"column:id;primary_key;autoIncrement" json:"id"` // ID
	BizId       string    `gorm:"column:biz_id" json:"bizId"`                    // 关联ID
	TenantId    string    `gorm:"column:tenant_id" json:"tenantId"`              // 租户ID
	Name        string    `gorm:"column:name" json:"name"`                       // 联系人名字
	State       uint8     `gorm:"column:state" json:"state"`                     // 状态 0:停用 1:启用
	Description string    `gorm:"column:description" json:"description"`         // 描述
	CreateUser  string    `gorm:"column:create_user" json:"createUser"`          // 创建人
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`          // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`          // 修改时间
}

func (*Contact) TableName() string {
	return "t_contact"
}
