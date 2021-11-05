package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/global"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertContactGroupCtl struct {
	dao *dao.AlertContactGroupDao
}

func NewAlertContactGroupCtl(dao *dao.AlertContactGroupDao) *AlertContactGroupCtl {
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
	var param forms.AlertContactGroupParam
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
	var param forms.AlertContactGroupParam
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
	var param forms.AlertContactGroupParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	mpc.dao.DeleteAlertContactGroup(param)
	c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
}
