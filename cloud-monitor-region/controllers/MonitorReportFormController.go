package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorReportFormCtl struct {
	service *service.MonitorReportFormService
}

func NewMonitorReportFormController(service *service.MonitorReportFormService) *MonitorReportFormCtl {
	return &MonitorReportFormCtl{service}
}

func (mpc *MonitorReportFormCtl) GetData(c *gin.Context) {
	var param forms.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	param.Labels = "instance,device"
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.service.GetData(param)))
}
