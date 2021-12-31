package model

type AlertContactInformation struct {
	Id         string `orm:"id" json:"id"`                   // ID
	TenantId   string `orm:"tenant_id" json:"tenant_id"`     // 租户ID
	ContactId  string `orm:"contact_id" json:"contact_id"`   // 联系人ID
	No         string `orm:"no" json:"no"`                   // 号码
	Type       int    `orm:"status" json:"status"`           // 类型 1:电话 2:邮箱 3:蓝信
	IsCertify  int    `orm:"is_certify" json:"is_certify"`   // 是否认证 1:已认证 2:未认证
	ActiveCode string `orm:"active_code" json:"active_code"` // 验证码
	CreateUser string `orm:"create_user" json:"create_user"` // 创建人
	CreateTime string `orm:"create_time" json:"create_time"` // 创建时间
	UpdateTime string `orm:"update_time" json:"update_time"` // 修改时间
}

func (*AlertContactInformation) TableName() string {
	return "alert_contact_information"
}
