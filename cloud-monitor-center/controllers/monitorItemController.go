package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorItemCtl struct {
	dao *dao.MonitorItemDao
}

func NewMonitorItemCtl(dao *dao.MonitorItemDao) *MonitorItemCtl {
	return &MonitorItemCtl{dao}
}

func (mpc *MonitorItemCtl) GetMonitorItemsById(c *gin.Context) {
	productId := c.Query("productId")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.SelectMonitorItemsById(productId)))
}
