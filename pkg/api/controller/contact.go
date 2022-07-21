package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	service2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContactCtl struct {
	service *service.ContactService
}

func NewContactCtl() *ContactCtl {
	return &ContactCtl{service.NewContactService(service.NewContactGroupService(service.NewContactGroupRelService()),
		service.NewContactInformationService(service2.NewMessageService(message_center.NewService())), service.NewContactGroupRelService())}
}

func (ctl *ContactCtl) GetContact(c *gin.Context) {
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
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.service.SelectContact(param)))
}

func (ctl *ContactCtl) CreateContact(c *gin.Context) {
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
	err = ctl.service.Persistence(ctl.service, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.Set(global.ResourceName, param.ContactBizId)
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (ctl *ContactCtl) UpdateContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.ContactBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enum.UpdateContact
	err = ctl.service.Persistence(ctl.service, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (ctl *ContactCtl) DeleteContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.ContactBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	param.EventEum = enum.DeleteContact
	err = ctl.service.Persistence(ctl.service, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}

func (ctl *ContactCtl) ActivateContact(c *gin.Context) {
	var param form.ContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.ActiveCode)
	param.EventEum = enum.ActivateContact
	err = ctl.service.Persistence(ctl.service, &param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("激活成功", ctl.getTenantName(param.ActiveCode)))
	}
}

//获取租户姓名
func (ctl *ContactCtl) getTenantName(activeCode string) string {
	tenantName := service2.NewTenantService().GetTenantInfo(ctl.service.GetTenantId(activeCode)).Name
	if strutil.IsBlank(tenantName) {
		return "未命名"
	}
	return tenantName
}
