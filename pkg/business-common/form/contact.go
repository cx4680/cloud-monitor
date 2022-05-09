package form

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"time"
)

type ContactForm struct {
	ContactBizId string    `gorm:"column:contact_biz_id" json:"contactBizId"`
	ContactName  string    `gorm:"column:contact_name" json:"contactName"`
	GroupBizId   string    `gorm:"column:group_biz_id" json:"groupBizId"`
	GroupName    string    `gorm:"column:group_name" json:"groupName"`
	Phone        string    `gorm:"column:phone" json:"phone"`
	PhoneState   int       `gorm:"column:phone_state" json:"phoneState"`
	Email        string    `gorm:"column:email" json:"email"`
	EmailState   int       `gorm:"column:email_state" json:"emailState"`
	Lanxin       string    `gorm:"column:lanxin" json:"lanxin"`
	LanxinState  int       `gorm:"column:lanxin_state" json:"lanxinState"`
	Description  string    `gorm:"column:description" json:"description"`
	GroupCount   int       `gorm:"column:group_count" json:"groupCount"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

type ContactGroupForm struct {
	GroupBizId   string    `gorm:"column:group_biz_id" json:"groupBizId"`
	GroupName    string    `gorm:"column:group_name" json:"groupName"`
	Description  string    `gorm:"column:description" json:"description"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
	ContactCount int       `gorm:"column:contact_count" json:"contactCount"`
}

type ContactFormPage struct {
	Records interface{} `json:"records"`
	Current int         `json:"current"`
	Size    int         `json:"size"`
	Total   int64       `json:"total"`
}

type ContactParam struct {
	TenantId         string   `form:"tenantId"`
	ContactBizId     string   `form:"contactBizId"`
	ContactName      string   `form:"contactName"`
	GroupBizId       string   `form:"groupBizId"`
	GroupName        string   `form:"groupName"`
	Phone            string   `form:"phone"`
	Email            string   `form:"email"`
	CreateUser       string   `form:"createUser"`
	Description      string   `form:"description"`
	ActiveCode       string   `form:"activeCode"`
	PageCurrent      int      `form:"pageCurrent"`
	PageSize         int      `form:"pageSize"`
	ContactBizIdList []string `form:"contactBizIdList"`
	GroupBizIdList   []string `form:"groupBizIdList"`
	EventEum         enum.EventEum
}
