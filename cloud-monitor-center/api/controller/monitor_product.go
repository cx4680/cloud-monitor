package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorProductCtl struct {
	service service.MonitorProductService
}

func NewMonitorProductCtl(service service.MonitorProductService) *MonitorProductCtl {
	return &MonitorProductCtl{service}
}

var MonitorProductService = service.NewMonitorProductService(dao.MonitorProduct)

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
	param.EventEum = enum.ChangeMonitorProductStatus
	err = mpc.service.Persistence(MonitorProductService, sys_rocketmq.MonitorProductTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}
