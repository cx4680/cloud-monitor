package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
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
	service.NewAlertContactInformationService(commonService.NewMessageService(messageCenter.NewService())), service.NewAlertContactGroupRelService())

func (acl *AlertContactCtl) GetAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acl.service.SelectAlertContact(param)))
}

func (acl *AlertContactCtl) InsertAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enums.InsertAlertContact
	err = alertContactService.Persistence(alertContactService, sysRocketMq.AlertContactTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (acl *AlertContactCtl) UpdateAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enums.UpdateAlertContact
	err = alertContactService.Persistence(alertContactService, sysRocketMq.AlertContactTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (acl *AlertContactCtl) DeleteAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enums.DeleteAlertContact
	err = alertContactService.Persistence(alertContactService, sysRocketMq.AlertContactTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}

func (acl *AlertContactCtl) CertifyAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	param.EventEum = enums.CertifyAlertContact
	err = alertContactService.Persistence(alertContactService, sysRocketMq.AlertContactTopic, param)
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
