package models

type AlertContact struct {
	Id          string `orm:"id" json:"id"`                   // ID
	TenantId    string `orm:"tenant_id" json:"tenantId"`      // 租户ID
	Name        string `orm:"name" json:"name"`               // 联系人名字
	Status      int    `orm:"status" json:"status"`           // 状态 0:停用 1:启用
	Description string `orm:"description" json:"description"` // 描述
	CreateUser  string `orm:"create_user" json:"createUser"`  // 创建人
	CreateTime  string `orm:"create_time" json:"createTime"`  // 创建时间
	UpdateTime  string `orm:"update_time" json:"updateTime"`  // 修改时间
}

func (*AlertContact) TableName() string {
	return "alert_contact"
}

type AlertContactGroup struct {
	Id          string `orm:"id" json:"id"`                   // ID
	TenantId    string `orm:"tenant_id" json:"tenant_id"`     // 租户ID
	Name        string `orm:"name" json:"name"`               // 组名
	Description string `orm:"description" json:"description"` // 描述
	CreateUser  string `orm:"create_user" json:"create_user"` // 创建人
	CreateTime  string `orm:"create_time" json:"create_time"` // 创建时间
	UpdateTime  string `orm:"update_time" json:"update_time"` // 修改时间
}

func (*AlertContactGroup) TableName() string {
	return "alert_contact_group"
}

type AlertContactGroupRel struct {
	Id         string `orm:"id" json:"id"`                   // ID
	TenantId   string `orm:"tenant_id" json:"tenant_id"`     // 租户ID
	ContactId  string `orm:"contact_id" json:"contact_id"`   // 联系人ID
	GroupId    string `orm:"group_id" json:"group_id"`       // 组ID
	CreateUser string `orm:"create_user" json:"create_user"` // 创建人
	CreateTime string `orm:"create_time" json:"create_time"` // 创建时间
	UpdateTime string `orm:"update_time" json:"update_time"` // 修改时间
}

func (*AlertContactGroupRel) TableName() string {
	return "alert_contact_group_rel"
}

type AlertContactInformation struct {
	Id         string `orm:"id" json:"id"`                   // ID
	TenantId   string `orm:"tenant_id" json:"tenant_id"`     // 租户ID
	ContactId  string `orm:"contact_id" json:"contact_id"`   // 联系人ID
	No         string `orm:"no" json:"no"`                   // 号码
	Type       int    `orm:"status" json:"status"`           // 类型 1:电话 2:邮箱 3:蓝信
	IsCertify  int    `orm:"is_certify" json:"is_certify"`   // 是否认证
	ActiveCode string `orm:"active_code" json:"active_code"` // 验证码
	CreateUser string `orm:"create_user" json:"create_user"` // 创建人
	CreateTime string `orm:"create_time" json:"create_time"` // 创建时间
	UpdateTime string `orm:"update_time" json:"update_time"` // 修改时间
}

func (*AlertContactInformation) TableName() string {
	return "alert_contact_information"
}
