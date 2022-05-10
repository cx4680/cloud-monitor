package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	service2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ContactCtl struct {
	service *service.ContactService
}

func NewContactCtl() *ContactCtl {
	return &ContactCtl{service.NewContactService(service.NewContactGroupService(service.NewContactGroupRelService()),
		service.NewContactInformationService(service2.NewMessageService(message_center.NewService())), service.NewContactGroupRelService())}
}

func (acl *ContactCtl) SelectContactPage(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	var param = ContactParam{PageNumber: 1, PageSize: 10}
	err = c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	result := acl.service.SelectContact(form.ContactParam{
		TenantId:    tenantId,
		ContactName: param.ContactName,
		Phone:       param.Phone,
		Email:       param.Mail,
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

func (acl *ContactCtl) CreateContact(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	var param = ContactParam{}
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.InvalidParameter, c))
		return
	}
	request := &form.ContactParam{
		TenantId:       tenantId,
		ContactName:    param.ContactName,
		Phone:          param.Phone,
		Email:          param.Mail,
		Description:    param.Description,
		CreateUser:     service2.NewTenantService().GetTenantInfo(tenantId).Name,
		GroupBizIdList: param.GroupIdList,
		EventEum:       enum.InsertContact,
	}
	err = acl.service.Persistence(acl.service, sys_rocketmq.ContactTopic, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(getErrorCode(err), c))
		return
	} else {
		result := struct {
			RequestId string
			ContactId string
		}{RequestId: openapi.GetRequestId(c), ContactId: request.ContactBizId}
		c.Set(global.ResourceName, request.ContactBizId)
		c.JSON(http.StatusOK, result)
	}
}

func (acl *ContactCtl) UpdateContact(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	var param = ContactParam{}
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.InvalidParameter, c))
		return
	}
	request := &form.ContactParam{
		TenantId:       tenantId,
		ContactBizId:   c.Param("ContactId"),
		ContactName:    param.ContactName,
		Phone:          param.Phone,
		Email:          param.Mail,
		Description:    param.Description,
		CreateUser:     service2.NewTenantService().GetTenantInfo(tenantId).Name,
		GroupBizIdList: param.GroupIdList,
		EventEum:       enum.UpdateContact,
	}
	c.Set(global.ResourceName, request.ContactBizId)
	err = acl.service.Persistence(acl.service, sys_rocketmq.ContactTopic, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(getErrorCode(err), c))
		return
	} else {
		result := struct{ RequestId string }{RequestId: openapi.GetRequestId(c)}
		c.JSON(http.StatusOK, result)
	}
}

func (acl *ContactCtl) DeleteContact(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	request := &form.ContactParam{
		TenantId:     tenantId,
		ContactBizId: c.Param("ContactId"),
		EventEum:     enum.DeleteContact,
	}
	c.Set(global.ResourceName, request.ContactBizId)
	err = acl.service.Persistence(acl.service, sys_rocketmq.ContactTopic, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(getErrorCode(err), c))
		return
	} else {
		result := struct{ RequestId string }{RequestId: openapi.GetRequestId(c)}
		c.JSON(http.StatusOK, result)
	}
}

func (acl *ContactCtl) ActivateContact(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	request := &form.ContactParam{
		TenantId:   tenantId,
		ActiveCode: c.Param("ActiveCode"),
		EventEum:   enum.ActivateContact,
	}
	c.Set(global.ResourceName, request.ActiveCode)
	err = acl.service.Persistence(acl.service, sys_rocketmq.ContactTopic, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(getErrorCode(err), c))
		return
	} else {
		result := struct {
			RequestId  string
			TenantName string
		}{
			RequestId:  openapi.GetRequestId(c),
			TenantName: service2.NewTenantService().GetTenantInfo(tenantId).Name,
		}
		c.JSON(http.StatusOK, result)
	}
}

func getErrorCode(err error) *openapi.ErrorCode {
	return ContactErrorMap[err.(*errors.BusinessError).Message]
}

type ContactParam struct {
	ContactName string   `form:"ContactName"`
	Phone       string   `form:"Phone"`
	Mail        string   `form:"Mail"`
	Description string   `form:"Description"`
	ActiveCode  string   `form:"ActiveCode"`
	PageNumber  int      `form:"PageNumber"`
	PageSize    int      `form:"PageSize"`
	GroupIdList []string `form:"GroupIdList"`
}

type ContactPage struct {
	RequestId  string    `json:"RequestId"`
	TotalCount int64     `json:"TotalCount"`
	PageSize   int       `json:"PageSize"`
	PageNumber int       `json:"PageNumber"`
	Contacts   []Contact `json:"Contacts"`
}

type Contact struct {
	ContactId   string    `json:"ContactId"`
	ContactName string    `json:"ContactName"`
	GroupNames  string    `json:"GroupNames"`
	Channels    []Channel `json:"Channels"`
	CreateTime  string    `json:"CreateTime"`
	UpdateTime  string    `json:"UpdateTime"`
	Description string    `json:"Description"`
}

type Channel struct {
	Channel string `json:"Channel"`
	Address string `json:"Address"`
	Status  int    `json:"Status"`
}

var ContactErrorMap = map[string]*openapi.ErrorCode{
	"手机号格式错误":      openapi.PhoneFormatError,
	"邮箱格式错误":       openapi.MailFormatError,
	"手机号和邮箱必须填写一项": openapi.InformationMissing,
	"联系人名字不能为空":    openapi.MissingContactName,
	"联系人ID不能为空":    openapi.MissingContactId,
	"联系组名字不能为空":    openapi.MissingGroupName,
	"联系组ID不能为空":    openapi.MissingGroupId,
	"联系人名字格式错误":    openapi.ContactNameFormatError,
	"联系组名字格式错误":    openapi.GroupNameFormatError,
	"请至少选择一位联系人":   openapi.MissingContact,
	"该租户无此联系人":     openapi.TenantHaveNotContact,
	"该租户无此联系组":     openapi.TenantHaveNotGroup,
	"联系组名重复":       openapi.GroupNameRepeat,
	"无效激活码":        openapi.InvalidActivationCode,
	"每个联系人最多加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组": openapi.ContactGroupNumberExceeded,
	"有联系人已加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组":   openapi.ContactGroupNumberExceeded,
	"有联系组已有" + strconv.Itoa(constant.MaxContactNum) + "个联系人":      openapi.GroupContactNumberExceeded,
	"联系人限制创建" + strconv.Itoa(constant.MaxContactNum) + "个":        openapi.ContactNumberExceeded,
	"联系组限制创建" + strconv.Itoa(constant.MaxGroupNum) + "个":          openapi.GroupNumberExceeded,
}
