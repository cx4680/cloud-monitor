package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContactCtl struct {
	service service.ContactService
}

func NewContactCtl(service service.ContactService) *ContactCtl {
	return &ContactCtl{service}
}

var contactService = service.NewContactService(service.NewContactGroupService(service.NewContactGroupRelService()),
	service.NewContactInformationService(commonService.NewMessageService(message_center.NewService())), service.NewContactGroupRelService())

func (acl *ContactCtl) GetContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acl.service.SelectContact(param)))
}

func (acl *ContactCtl) AddContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enum.InsertContact
	err = contactService.Persistence(contactService, sys_rocketmq.ContactTopic, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (acl *ContactCtl) UpdateContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enum.UpdateContact
	err = contactService.Persistence(contactService, sys_rocketmq.ContactTopic, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (acl *ContactCtl) DeleteContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enum.DeleteContact
	err = contactService.Persistence(contactService, sys_rocketmq.ContactTopic, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}

func (acl *ContactCtl) ActivateContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	param.EventEum = enum.ActivateContact
	err = contactService.Persistence(contactService, sys_rocketmq.ContactTopic, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("激活成功", getTenantName(param.ActiveCode)))
	}
}

//获取租户姓名
func getTenantName(activeCode string) string {
	tenantName := commonService.NewTenantService().GetTenantInfo(contactService.GetTenantId(activeCode)).Name
	if strutil.IsBlank(tenantName) {
		return "未命名"
	}
	return tenantName
}
