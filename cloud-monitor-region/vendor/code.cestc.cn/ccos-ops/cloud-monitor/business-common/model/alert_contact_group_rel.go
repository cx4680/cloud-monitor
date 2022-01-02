package model

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"time"
)

type AlertContactGroupRel struct {
	Id         string    `orm:"id;primary_key" json:"id"`       // ID
	TenantId   string    `orm:"tenant_id" json:"tenant_id"`     // 租户ID
	ContactId  string    `orm:"contact_id" json:"contact_id"`   // 联系人ID
	GroupId    string    `orm:"group_id" json:"group_id"`       // 组ID
	CreateUser string    `orm:"create_user" json:"create_user"` // 创建人
	CreateTime time.Time `orm:"create_time" json:"create_time"` // 创建时间
	UpdateTime time.Time `orm:"update_time" json:"update_time"` // 修改时间
}

func (*AlertContactGroupRel) TableName() string {
	return "alert_contact_group_rel"
}

type UpdateAlertContactGroupRel struct {
	RelList []*AlertContactGroupRel
	Param   form.AlertContactParam
}
