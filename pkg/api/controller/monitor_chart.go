package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
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

func (ctl *MonitorChartCtl) GetData(c *gin.Context) {
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
	data, err := ctl.service.GetData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (ctl *MonitorChartCtl) GetAxisData(c *gin.Context) {
	var param = form.PrometheusRequest{Step: 60, Scope: "1m"}
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
	data, err := ctl.service.GetAxisData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (ctl *MonitorChartCtl) GetTopData(c *gin.Context) {
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
	data, err := ctl.service.GetTop(*param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (ctl *MonitorChartCtl) GetNetworkData(c *gin.Context) {
	var param form.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	data, err := ctl.service.GetNetworkData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (ctl *MonitorChartCtl) GetAxisDataInner(c *gin.Context) {
	var param = form.PrometheusRequest{Step: 60, Scope: "1m"}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	data, err := ctl.service.GetAxisData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (ctl *MonitorChartCtl) GetProcessData(c *gin.Context) {
	var param = form.PrometheusRequest{Step: 60}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	data, err := ctl.service.GetProcessData(param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (ctl *MonitorChartCtl) GetTopDataByIam(c *gin.Context) {
	tenantId, iamUserId, err := util.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	if strutil.IsBlank(iamUserId) || iamUserId == tenantId {
		ctl.GetTopData(c)
		return
	}
	var param = &form.PrometheusRequest{TopNum: 5}
	err = c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.Name)

	param.TenantId = tenantId
	param.IamUserId = iamUserId
	data, err := ctl.service.GetTopDataByIam(*param)
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}
