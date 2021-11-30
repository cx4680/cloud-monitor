package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertContactCtl struct {
	service service.AlertContactService
}

func NewAlertContactCtl(service service.AlertContactService) *AlertContactCtl {
	return &AlertContactCtl{service}
}

var groupService = service.NewAlertContactGroupService()
var informationService = service.NewAlertContactInformationService()
var contactService = service.NewAlertContactService(groupService, informationService)

func (acl *AlertContactCtl) GetAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", acl.service.Select(param)))
}

func (acl *AlertContactCtl) InsertAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	param.EventEum = enums.InsertAlertContact
	err = contactService.Persistence(contactService, sysRocketMq.AlertContactTopic, param)
	//local, err := acl.service.PersistenceLocal(global.DB, param)
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
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	err = acl.service.Update(param)
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
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	err = acl.service.Delete(param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
	}
}

func (acl *AlertContactCtl) CertifyAlertContact(c *gin.Context) {
	activeCode := c.Query("activeCode")
	c.JSON(http.StatusOK, global.NewSuccess("激活成功", acl.service.CertifyAlertContact(activeCode)))
}
