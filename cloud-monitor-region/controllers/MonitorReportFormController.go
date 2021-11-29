package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/ecs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

func (mpc *MonitorReportFormCtl) GetAxisData(c *gin.Context) {
	var param forms.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	monitorItemDao := dao.MonitorItem
	param.Labels = monitorItemDao.GetLabelsByName(param.Name)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.service.GetAxisData(param)))
}

func (mpc *MonitorReportFormCtl) GetTop(c *gin.Context) {
	var param forms.PrometheusRequest
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	var form = forms.EcsQueryPageForm{
		TenantId: param.TenantId,
		Current:  1,
		PageSize: 1000,
	}
	rows, err := ecs.PageList(&form)
	var instanceList []string
	for _, ecsVO := range rows.Rows {
		if ecsVO.InstanceId == "" {
			instanceList = append(instanceList, ecsVO.InstanceId)
		}
	}
	param.Instance = strings.Join(instanceList, "|")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.service.GetTop(param)))
}
