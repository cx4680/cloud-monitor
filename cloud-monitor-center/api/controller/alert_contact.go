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
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertContactCtl struct {
	service service.AlertContactService
}

func NewAlertContactCtl(service service.AlertContactService) *AlertContactCtl {
	return &AlertContactCtl{service}
}

var alertContactService = service.NewAlertContactService(service.NewAlertContactGroupService(service.NewAlertContactGroupRelService()),
	service.NewAlertContactInformationService(commonService.NewMessageService(message_center.NewService())), service.NewAlertContactGroupRelService())

func (acl *AlertContactCtl) GetAlertContact(c *gin.Context) {
	var param form.AlertContactParam
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
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acl.service.SelectAlertContact(param)))
}

func (acl *AlertContactCtl) InsertAlertContact(c *gin.Context) {
	var param form.AlertContactParam
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
	param.EventEum = enum.InsertAlertContact
	err = alertContactService.Persistence(alertContactService, sys_rocketmq.AlertContactTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (acl *AlertContactCtl) UpdateAlertContact(c *gin.Context) {
	var param form.AlertContactParam
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
	param.EventEum = enum.UpdateAlertContact
	err = alertContactService.Persistence(alertContactService, sys_rocketmq.AlertContactTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (acl *AlertContactCtl) DeleteAlertContact(c *gin.Context) {
	var param form.AlertContactParam
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
	param.EventEum = enum.DeleteAlertContact
	err = alertContactService.Persistence(alertContactService, sys_rocketmq.AlertContactTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}

func (acl *AlertContactCtl) CertifyAlertContact(c *gin.Context) {
	var param form.AlertContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	param.EventEum = enum.CertifyAlertContact
	err = alertContactService.Persistence(alertContactService, sys_rocketmq.AlertContactTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("激活成功", getTenantName(param.ActiveCode)))
	}
}

//获取租户姓名
func getTenantName(activeCode string) string {
	tenantName := commonService.NewTenantService().GetTenantInfo(alertContactService.GetTenantId(activeCode)).Name
	if tenantName == "" {
		return "未命名"
	}
	return tenantName
}
