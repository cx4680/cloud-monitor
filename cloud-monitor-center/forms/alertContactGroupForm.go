package forms

type AlertContactGroupForm struct {
	GroupId     string `orm:"group_id" json:"groupId"`
	GroupName   string `orm:"group_name" json:"groupName"`
	Description string `orm:"description" json:"description"`
	CreateTime  string `orm:"create_time" json:"createTime"`
	UpdateTime  string `orm:"update_time" json:"updateTime"`
}

type AlertContactGroupParam struct {
	TenantId      string
	GroupId       string
	GroupName     string
	Description   string
	CreateUser    string
	ContactIdList []string
	PageCurrent   string
	PageSize      string
}
