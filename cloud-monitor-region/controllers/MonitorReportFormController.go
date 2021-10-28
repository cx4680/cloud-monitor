package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/ecs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
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
	monitorItemDao := dao.NewMonitorItemDao(database.GetDb())
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
	pageResult, _ := ecs.GetUserInstancePage(nil, 0, 100, param.TenantId)
	rows := pageResult.Rows
	var instanceList = []string{"ecs-ce2qxixq28fu0q", "ecs-ce310hziqnk67l", "ecs-ce3293m18eimi7"}
	for i := range rows {
		if rows[i].HostId == "" {
			instanceList = append(instanceList, rows[i].HostId)
		}
	}
	param.Instance = strings.Join(instanceList, "|")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.service.GetTop(param)))
}
