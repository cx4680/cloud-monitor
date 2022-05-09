package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContactGroupCtl struct {
	service *service.ContactGroupService
}

func NewContactGroupCtl() *ContactGroupCtl {
	return &ContactGroupCtl{service.NewContactGroupService(service.NewContactGroupRelService())}
}

func (ctl *ContactGroupCtl) GetContactGroup(c *gin.Context) {
	var param = form.ContactParam{PageCurrent: 1, PageSize: 10}
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
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.service.SelectContactGroup(param)))
}

func (ctl *ContactGroupCtl) GetGroupContact(c *gin.Context) {
	var param = form.ContactParam{PageCurrent: 1, PageSize: 10}
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
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.service.SelectAlertGroupContact(param)))
}

func (ctl *ContactGroupCtl) CreateContactGroup(c *gin.Context) {
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
	err = ctl.service.Persistence(ctl.service, sys_rocketmq.ContactGroupTopic, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.Set(global.ResourceName, param.GroupBizId)
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (ctl *ContactGroupCtl) UpdateContactGroup(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.GroupBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enum.UpdateContactGroup
	err = ctl.service.Persistence(ctl.service, sys_rocketmq.ContactGroupTopic, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (ctl *ContactGroupCtl) DeleteContactGroup(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.GroupBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enum.DeleteContactGroup
	err = ctl.service.Persistence(ctl.service, sys_rocketmq.ContactGroupTopic, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}

func (ctl *ContactGroupCtl) GetContactGroupWithSys(c *gin.Context) {
	var param = form.ContactParam{PageCurrent: 1, PageSize: 10000}
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
	groupPage := ctl.service.SelectContactGroup(param)
	if len(groupPage.Records.([]form.ContactGroupForm)) == 0 {
		groupPage.Records = append(groupPage.Records.([]form.ContactGroupForm), form.ContactGroupForm{
			GroupBizId:   "-1",
			GroupName:    constant.DefaultContact,
			Description:  "系统创建",
			CreateTime:   util.GetNow(),
			UpdateTime:   util.GetNow(),
			ContactCount: 1,
		})
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", groupPage))
}
