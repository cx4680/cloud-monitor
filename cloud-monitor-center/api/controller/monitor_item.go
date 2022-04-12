package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type MonitorItemCtl struct {
	service service.MonitorItemService
}

func NewMonitorItemCtl(service service.MonitorItemService) *MonitorItemCtl {
	return &MonitorItemCtl{service}
}

var MonitorItemService = service.NewMonitorItemService(dao.MonitorItem)

var displayList = []string{"chart", "rule", "scaling"}

func (mic *MonitorItemCtl) GetMonitorItemsById(c *gin.Context) {
	var param form.MonitorItemParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.ProductBizId)
	if strutil.IsNotBlank(param.Display) && !checkDisplay(param.Display) {
		c.JSON(http.StatusOK, global.NewError("查询失败，展示参数错误"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mic.service.GetMonitorItem(param)))
}

func (mic *MonitorItemCtl) ChangeDisplay(c *gin.Context) {
	var param form.MonitorItemParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.ProductBizId)
	for _, v := range strings.Split(param.Display, ",") {
		if !checkDisplay(v) {
			c.JSON(http.StatusOK, global.NewError("修改失败，展示参数错误"))
			return
		}
	}
	param.EventEum = enum.ChangeMonitorItemDisplay
	err = mic.service.Persistence(MonitorItemService, sys_rocketmq.MonitorItemTopic, param)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("修改成功", true))
	}
}

func checkDisplay(display string) bool {
	for _, v := range displayList {
		if display == v {
			return true
		}
	}
	return false
}
