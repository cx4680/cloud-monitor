package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/global"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertContactCtl struct {
	dao *dao.AlertContactDao
}

func NewAlertContactCtl(dao *dao.AlertContactDao) *AlertContactCtl {
	return &AlertContactCtl{dao}
}

func (mpc *AlertContactCtl) GetAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.GetAlertContact(param)))
}

func (mpc *AlertContactCtl) InsertAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	err = mpc.dao.InsertAlertContact(param)
	if err.Error() != "" {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
	}
}

func (mpc *AlertContactCtl) UpdateAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	mpc.dao.UpdateAlertContact(param)
	c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
}

func (mpc *AlertContactCtl) DeleteAlertContact(c *gin.Context) {
	var param forms.AlertContactParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	mpc.dao.DeleteAlertContact(param)
	c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
}

func (mpc *AlertContactCtl) CertifyAlertContact(c *gin.Context) {
	activeCode := c.Query("activeCode")
	c.JSON(http.StatusOK, global.NewSuccess("激活成功", mpc.dao.CertifyAlertContact(activeCode)))
}
