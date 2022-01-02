package model

import "time"

type AlertContactGroup struct {
	Id          string    `orm:"id;primary_key" json:"id"`       // ID
	TenantId    string    `orm:"tenant_id" json:"tenant_id"`     // 租户ID
	Name        string    `orm:"name" json:"name"`               // 组名
	Description string    `orm:"description" json:"description"` // 描述
	CreateUser  string    `orm:"create_user" json:"create_user"` // 创建人
	CreateTime  time.Time `orm:"create_time" json:"create_time"` // 创建时间
	UpdateTime  time.Time `orm:"update_time" json:"update_time"` // 修改时间
}

func (*AlertContactGroup) TableName() string {
	return "alert_contact_group"
}
