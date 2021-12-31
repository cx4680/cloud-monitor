package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertContactGroupCtl struct {
	service service.AlertContactGroupService
}

func NewAlertContactGroupCtl(service service.AlertContactGroupService) *AlertContactGroupCtl {
	return &AlertContactGroupCtl{service}
}

var alertContactGroupService = service.NewAlertContactGroupService(service.NewAlertContactGroupRelService())

func (acgc *AlertContactGroupCtl) GetAlertContactGroup(c *gin.Context) {
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
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acgc.service.SelectAlertContactGroup(param)))
}

func (acgc *AlertContactGroupCtl) GetAlertGroupContact(c *gin.Context) {
	var param form.AlertContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acgc.service.SelectAlertGroupContact(param)))
}

func (acgc *AlertContactGroupCtl) InsertAlertContactGroup(c *gin.Context) {
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
	param.EventEum = enum.InsertAlertContactGroup
	err = alertContactGroupService.Persistence(alertContactGroupService, sys_rocketmq.AlertContactGroupTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (acgc *AlertContactGroupCtl) UpdateAlertContactGroup(c *gin.Context) {
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
	param.EventEum = enum.UpdateAlertContactGroup
	err = alertContactGroupService.Persistence(alertContactGroupService, sys_rocketmq.AlertContactGroupTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (acgc *AlertContactGroupCtl) DeleteAlertContactGroup(c *gin.Context) {
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
	param.EventEum = enum.DeleteAlertContactGroup
	err = alertContactGroupService.Persistence(alertContactGroupService, sys_rocketmq.AlertContactGroupTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}
