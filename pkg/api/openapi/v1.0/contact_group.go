package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContactGroupCtl struct {
	service *service.ContactGroupService
}

func NewContactGroupCtl() *ContactGroupCtl {
	return &ContactGroupCtl{service.NewContactGroupService(service.NewContactGroupRelService())}
}

func (acgc *ContactGroupCtl) SelectContactGroupPage(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	var param = ContactGroupParam{PageNumber: 1, PageSize: 10}
	err = c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	result := acgc.service.SelectContactGroup(form.ContactParam{
		TenantId:    tenantId,
		GroupName:   param.GroupName,
		PageCurrent: param.PageNumber,
		PageSize:    param.PageSize,
	})
	var contactGroups []ContactGroup
	for _, v := range result.Records.([]form.ContactGroupForm) {
		contactGroup := ContactGroup{
			GroupId:      v.GroupBizId,
			GroupName:    v.GroupName,
			ContactCount: v.ContactCount,
			CreateTime:   util.TimeToStr(v.CreateTime, util.FullTimeFmt),
			UpdateTime:   util.TimeToStr(v.UpdateTime, util.FullTimeFmt),
			Description:  v.Description,
		}
		contactGroups = append(contactGroups, contactGroup)
	}
	contactPage := ContactGroupPage{
		RequestId:     openapi.GetRequestId(c),
		PageNumber:    result.Current,
		PageSize:      result.Size,
		TotalCount:    result.Total,
		ContactGroups: contactGroups,
	}
	c.JSON(http.StatusOK, contactPage)
}

func (acgc *ContactGroupCtl) SelectContactPageByGroupId(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	var param = ContactGroupParam{PageNumber: 1, PageSize: 10}
	result := acgc.service.SelectAlertGroupContact(form.ContactParam{
		TenantId:    tenantId,
		GroupBizId:  c.Param("GroupId"),
		PageCurrent: param.PageNumber,
		PageSize:    param.PageSize,
	})
	var contacts []Contact
	for _, v := range result.Records.([]form.ContactForm) {
		contact := Contact{
			ContactId:   v.ContactBizId,
			ContactName: v.ContactName,
			GroupNames:  v.GroupName,
			CreateTime:  util.TimeToStr(v.CreateTime, util.FullTimeFmt),
			UpdateTime:  util.TimeToStr(v.UpdateTime, util.FullTimeFmt),
			Description: v.Description,
		}
		if strutil.IsNotBlank(v.Phone) {
			channel := Channel{
				Channel: "phone",
				Address: v.Phone,
				Status:  v.PhoneState,
			}
			contact.Channels = append(contact.Channels, channel)
		}
		if strutil.IsNotBlank(v.Email) {
			channel := Channel{
				Channel: "mail",
				Address: v.Email,
				Status:  v.EmailState,
			}
			contact.Channels = append(contact.Channels, channel)
		}
		contacts = append(contacts, contact)
	}
	contactPage := ContactPage{
		RequestId:  openapi.GetRequestId(c),
		PageNumber: result.Current,
		PageSize:   result.Size,
		TotalCount: result.Total,
		Contacts:   contacts,
	}
	c.JSON(http.StatusOK, contactPage)
}

func (acgc *ContactGroupCtl) CreateContactGroup(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	var param = ContactGroupParam{}
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.InvalidParameter, c))
		return
	}
	request := &form.ContactParam{
		TenantId:         tenantId,
		GroupName:        param.GroupName,
		Description:      param.Description,
		CreateUser:       commonService.NewTenantService().GetTenantInfo(tenantId).Name,
		ContactBizIdList: param.ContactIdList,
		EventEum:         enum.InsertContactGroup,
	}
	err = acgc.service.Persistence(acgc.service, sys_rocketmq.ContactGroupTopic, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(ContactErrorMap[err.(*errors.BusinessError).Message], c))
		return
	} else {
		result := struct {
			RequestId string
			GroupId   string
		}{RequestId: openapi.GetRequestId(c), GroupId: request.GroupBizId}
		c.JSON(http.StatusOK, result)
	}
}

func (acgc *ContactGroupCtl) UpdateContactGroup(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	var param = ContactGroupParam{}
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.InvalidParameter, c))
		return
	}
	request := &form.ContactParam{
		TenantId:         tenantId,
		GroupBizId:       c.Param("GroupId"),
		GroupName:        param.GroupName,
		Description:      param.Description,
		CreateUser:       commonService.NewTenantService().GetTenantInfo(tenantId).Name,
		ContactBizIdList: param.ContactIdList,
		EventEum:         enum.UpdateContactGroup,
	}
	err = acgc.service.Persistence(acgc.service, sys_rocketmq.ContactGroupTopic, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(ContactErrorMap[err.(*errors.BusinessError).Message], c))
		return
	} else {
		result := struct{ RequestId string }{RequestId: openapi.GetRequestId(c)}
		c.JSON(http.StatusOK, result)
	}
}

func (acgc *ContactGroupCtl) DeleteContactGroup(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	request := &form.ContactParam{
		TenantId:   tenantId,
		GroupBizId: c.Param("GroupId"),
		EventEum:   enum.DeleteContactGroup,
	}
	err = acgc.service.Persistence(acgc.service, sys_rocketmq.ContactGroupTopic, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(ContactErrorMap[err.(*errors.BusinessError).Message], c))
		return
	} else {
		result := struct{ RequestId string }{RequestId: openapi.GetRequestId(c)}
		c.JSON(http.StatusOK, result)
	}
}

type ContactGroupParam struct {
	GroupName     string   `form:"GroupName"`
	Description   string   `form:"Description"`
	PageNumber    int      `form:"PageNumber"`
	PageSize      int      `form:"PageSize"`
	ContactIdList []string `form:"ContactIdList"`
}

type ContactGroupPage struct {
	RequestId     string         `json:"RequestId"`
	TotalCount    int64          `json:"TotalCount"`
	PageSize      int            `json:"PageSize"`
	PageNumber    int            `json:"PageNumber"`
	ContactGroups []ContactGroup `json:"Contacts"`
}

type ContactGroup struct {
	GroupId      string `json:"GroupId"`
	GroupName    string `json:"GroupName"`
	ContactCount int    `json:"ContactCount"`
	CreateTime   string `json:"CreateTime"`
	UpdateTime   string `json:"UpdateTime"`
	Description  string `json:"Description"`
}
