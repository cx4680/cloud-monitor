package form

import "code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"

type AlertContactForm struct {
	ContactId     string `orm:"contact_id" json:"contactId"`
	ContactName   string `orm:"contact_name" json:"contactName"`
	GroupId       string `orm:"group_id" json:"groupId"`
	GroupName     string `orm:"group_name" json:"groupName"`
	Phone         string `orm:"phone" json:"phone"`
	PhoneCertify  int    `orm:"phone_certify" json:"phoneCertify"`
	Email         string `orm:"email" json:"email"`
	EmailCertify  int    `orm:"email_certify" json:"emailCertify"`
	Lanxin        string `orm:"lanxin" json:"lanxin"`
	LanxinCertify int    `orm:"lanxin_certify" json:"lanxinCertify"`
	Description   string `orm:"description" json:"description"`
	GroupCount    int    `orm:"groupCount" json:"groupCount"`
}

type AlertContactGroupForm struct {
	GroupId      string `orm:"group_id" json:"groupId"`
	GroupName    string `orm:"group_name" json:"groupName"`
	Description  string `orm:"description" json:"description"`
	CreateTime   string `orm:"create_time" json:"createTime"`
	UpdateTime   string `orm:"update_time" json:"updateTime"`
	ContactCount int    `orm:"contactCount" json:"contactCount"`
}

type AlertContactFormPage struct {
	Records interface{} `json:"records"`
	Current int         `json:"current"`
	Size    int         `json:"size"`
	Total   int64       `json:"total"`
}

type AlertContactParam struct {
	TenantId      string   `form:"tenantId"`
	ContactId     string   `form:"contactId"`
	ContactName   string   `form:"contactName"`
	GroupId       string   `form:"groupId"`
	GroupName     string   `form:"groupName"`
	Phone         string   `form:"phone"`
	Email         string   `form:"email"`
	Lanxin        string   `form:"lanxin"`
	CreateUser    string   `form:"createUser"`
	Description   string   `form:"description"`
	ActiveCode    string   `form:"activeCode"`
	PageCurrent   int      `form:"pageCurrent"`
	PageSize      int      `form:"pageSize"`
	ContactIdList []string `form:"contactIdList"`
	GroupIdList   []string `form:"groupIdList"`
	EventEum      enum.EventEum
}
