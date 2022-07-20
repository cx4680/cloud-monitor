package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReportFormCtl struct {
	service *service.ReportFormService
}

func NewReportFormController() *ReportFormCtl {
	return &ReportFormCtl{
		service: service.NewReportFormService(),
	}
}

func (rfc *ReportFormCtl) GetMonitorData(c *gin.Context) {
	var callback = form.CallbackReportForm{}
	err := c.ShouldBindJSON(&callback)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	var param = form.ReportFormParam{Step: 60}
	jsonutil.ToObject(callback.Param, &param)
	param.RegionCode = config.Cfg.Common.RegionName
	if len(param.InstanceList) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("实例不能为空"))
		return
	}
	if len(param.ItemList) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("监控指标不能为空"))
		return
	}
	if param.Start == 0 || param.End == 0 || param.Start > param.End {
		c.JSON(http.StatusBadRequest, global.NewError("时间参数有误"))
		return
	}
	result, err := rfc.service.GetMonitorData(param)
	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":      http.StatusOK,
			"message":   "success",
			"pageCount": 1,
			"result":    result,
		})
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (rfc *ReportFormCtl) ExportMonitorData(c *gin.Context) {
	var param = form.ReportFormParam{Step: 60}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	if len(param.InstanceList) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("实例ID不能为空"))
		return
	}
	if len(param.ItemList) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("监控指标不能为空"))
		return
	}
	if param.Start == 0 || param.End == 0 || param.Start > param.End {
		c.JSON(http.StatusBadRequest, global.NewError("时间参数有误"))
		return
	}
	param.RegionCode = config.Cfg.Common.RegionName
	err = rfc.service.ExportMonitorData(param, c.Request.Header.Get("user-info"))
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("导入任务已下发", true))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (rfc *ReportFormCtl) GetAlarmRecord(c *gin.Context) {
	var callback = form.CallbackReportForm{}
	err := c.ShouldBindJSON(&callback)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	var param = form.AlarmRecordPageQueryForm{}
	jsonutil.ToObject(callback.Param, &param)
	result, err := rfc.service.GetAlarmRecord(param)
	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":      http.StatusOK,
			"message":   "success",
			"pageCount": 1,
			"result":    result,
		})
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}

func (rfc *ReportFormCtl) ExportAlarmRecord(c *gin.Context) {
	var param = form.AlarmRecordPageQueryForm{}
	err := c.ShouldBindJSON(&param)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.RegionCode = config.Cfg.Common.RegionName
	param.TenantID = tenantId
	err = rfc.service.ExportAlarmRecord(param, c.Request.Header.Get("user-info"))
	if err == nil {
		c.JSON(http.StatusOK, global.NewSuccess("导入任务已下发", true))
	} else {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	}
}
