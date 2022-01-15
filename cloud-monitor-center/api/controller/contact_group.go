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

type ContactGroupCtl struct {
	service service.ContactGroupService
}

func NewContactGroupCtl(service service.ContactGroupService) *ContactGroupCtl {
	return &ContactGroupCtl{service}
}

var contactGroupService = service.NewContactGroupService(service.NewContactGroupRelService())

func (acgc *ContactGroupCtl) GetContactGroup(c *gin.Context) {
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
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acgc.service.SelectContactGroup(param)))
}

func (acgc *ContactGroupCtl) GetGroupContact(c *gin.Context) {
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
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acgc.service.SelectAlertGroupContact(param)))
}

func (acgc *ContactGroupCtl) AddContactGroup(c *gin.Context) {
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
	param.EventEum = enum.InsertContactGroup
	err = contactGroupService.Persistence(contactGroupService, sys_rocketmq.ContactGroupTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (acgc *ContactGroupCtl) UpdateContactGroup(c *gin.Context) {
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
	param.EventEum = enum.UpdateContactGroup
	err = contactGroupService.Persistence(contactGroupService, sys_rocketmq.ContactGroupTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (acgc *ContactGroupCtl) DeleteContactGroup(c *gin.Context) {
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
	param.EventEum = enum.DeleteContactGroup
	err = contactGroupService.Persistence(contactGroupService, sys_rocketmq.ContactGroupTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}
