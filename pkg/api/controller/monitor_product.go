package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorProductCtl struct {
	service *service.MonitorProductService
}

func NewMonitorProductCtl() *MonitorProductCtl {
	return &MonitorProductCtl{service.NewMonitorProductService(dao.MonitorProduct)}
}

func (mpc *MonitorProductCtl) GetMonitorProduct(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.service.GetMonitorProduct()))
}

func (mpc *MonitorProductCtl) GetAllMonitorProduct(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.service.GetAllMonitorProduct()))
}

func (mpc *MonitorProductCtl) ChangeStatus(c *gin.Context) {
	var param form.MonitorProductParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.BizIdList)
	param.EventEum = enum.ChangeMonitorProductStatus
	err = mpc.service.Persistence(mpc.service, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}
