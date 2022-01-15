package model

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"time"
)

type ContactGroupRel struct {
	Id           uint64    `gorm:"column:id;primary_key;autoIncrement" json:"id"` // ID
	TenantId     string    `gorm:"column:tenant_id" json:"tenantId"`              // 租户ID
	ContactBizId string    `gorm:"column:contact_biz_id" json:"contactBizId"`     // 联系人ID
	GroupBizId   string    `gorm:"column:group_biz_id" json:"groupBizId"`         // 组ID
	CreateUser   string    `gorm:"column:create_user" json:"create_user"`         // 创建人
	CreateTime   time.Time `gorm:"column:create_time" json:"create_time"`         // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`         // 修改时间
}

func (*ContactGroupRel) TableName() string {
	return "t_contact_group_rel"
}

type UpdateContactGroupRel struct {
	RelList []*ContactGroupRel
	Param   form.ContactParam
}
