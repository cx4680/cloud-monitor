package model

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
