package model

import "time"

type ContactInformation struct {
	Id           uint64    `gorm:"column:id;primary_key;autoIncrement" json:"id"` // ID
	TenantId     string    `gorm:"column:tenant_id" json:"tenantId"`              // 租户ID
	ContactBizId string    `gorm:"column:contact_biz_id" json:"contactBizId"`     // 联系人ID
	Address      string    `gorm:"column:address" json:"address"`                 // 号码
	Type         uint8     `gorm:"column:type" json:"type"`                       // 类型 1:电话 2:邮箱 3:蓝信
	State        uint8     `gorm:"column:state" json:"state"`                     // 状态 1:已激活 0:未激活
	ActiveCode   string    `gorm:"column:active_code" json:"activeCode"`          // 验证码
	CreateUser   string    `gorm:"column:create_user" json:"createUser"`          // 创建人
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`          // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`          // 修改时间
}

func (*ContactInformation) TableName() string {
	return "t_contact_information"
}
