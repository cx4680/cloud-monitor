package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorChartCtl struct {
	service *service.MonitorChartService
}

func NewMonitorChartController(service *service.MonitorChartService) *MonitorChartCtl {
	return &MonitorChartCtl{service}
}

func (mcc *MonitorChartCtl) GetData(c *gin.Context) {
	var param form.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.Instance)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	data, err := mcc.service.GetData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (mcc *MonitorChartCtl) GetAxisData(c *gin.Context) {
	var param = form.PrometheusRequest{Step: 60}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.Instance)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	data, err := mcc.service.GetAxisData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (mcc *MonitorChartCtl) GetTop(c *gin.Context) {
	var param = &form.PrometheusRequest{TopNum: 5}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.Name)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	data, err := mcc.service.GetTop(*param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (mcc *MonitorChartCtl) GetNetworkData(c *gin.Context) {
	var param form.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	data, err := mcc.service.GetNetworkData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (mcc *MonitorChartCtl) GetAxisDataInner(c *gin.Context) {
	var param form.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	data, err := mcc.service.GetAxisData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}
