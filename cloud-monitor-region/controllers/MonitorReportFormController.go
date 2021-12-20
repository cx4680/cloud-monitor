package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
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
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	monitorItemDao := dao.MonitorItem
	param.Labels = monitorItemDao.GetMonitorItemByName(param.Name).Labels
	data, err := mpc.service.GetData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (mpc *MonitorReportFormCtl) GetAxisData(c *gin.Context) {
	var param forms.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	data, err := mpc.service.GetAxisData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (mpc *MonitorReportFormCtl) GetTop(c *gin.Context) {
	var param forms.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	data, err := mpc.service.GetTop(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}
