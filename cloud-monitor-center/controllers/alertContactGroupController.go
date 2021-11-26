package controllers

import (
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	forms2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertContactGroupCtl struct {
	dao *dao2.AlertContactGroupDao
}

func NewAlertContactGroupCtl(dao *dao2.AlertContactGroupDao) *AlertContactGroupCtl {
	return &AlertContactGroupCtl{dao}
}

func (mpc *AlertContactGroupCtl) GetAlertContactGroup(c *gin.Context) {
	tenantId := c.Query("tenantId")
	groupName := c.Query("groupName")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.GetAlertContactGroup(tenantId, groupName)))
}

func (mpc *AlertContactGroupCtl) GetAlertGroupContact(c *gin.Context) {
	tenantId := c.Query("tenantId")
	groupId := c.Query("groupId")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.GetAlertGroupContact(tenantId, groupId)))
}

func (mpc *AlertContactGroupCtl) InsertAlertContactGroup(c *gin.Context) {
	var param forms2.AlertContactGroupParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	err = mpc.dao.InsertAlertContactGroup(param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (mpc *AlertContactGroupCtl) UpdateAlertContactGroup(c *gin.Context) {
	var param forms2.AlertContactGroupParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	err = mpc.dao.UpdateAlertContactGroup(param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func (mpc *AlertContactGroupCtl) DeleteAlertContactGroup(c *gin.Context) {
	var param forms2.AlertContactGroupParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	mpc.dao.DeleteAlertContactGroup(param)
	c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
}
