package forms

import "code.cestc.cn/ccos-ops/cloud-monitor/common/enums"

type AlertContactForm struct {
	ContactId     string `orm:"contact_id" json:"contactId"`
	ContactName   string `orm:"contact_name" json:"contactName"`
	GroupId       string `orm:"group_id" json:"groupId"`
	GroupName     string `orm:"group_name" json:"groupName"`
	Phone         string `orm:"phone" json:"phone"`
	PhoneCertify  string `orm:"phone_certify" json:"phoneCertify"`
	Email         string `orm:"email" json:"email"`
	EmailCertify  string `orm:"email_certify" json:"emailCertify"`
	Lanxin        string `orm:"lanxin" json:"lanxin"`
	LanxinCertify string `orm:"lanxin_certify" json:"lanxinCertify"`
	Description   string `orm:"description" json:"description"`
}

type AlertContactFormPage struct {
	Records *[]AlertContactForm `json:"records"`
	Current int                 `json:"current"`
	Size    int                 `json:"size"`
	Total   int64               `json:"total"`
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
	EventEum      enums.EventEum
}

type MqMsg struct {
	EventEum enums.EventEum
	Data     interface{}
}
