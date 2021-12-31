package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorItemCtl struct {
	dao *dao.MonitorItemDao
}

func NewMonitorItemCtl(dao *dao.MonitorItemDao) *MonitorItemCtl {
	return &MonitorItemCtl{dao}
}

func (d *MonitorItemCtl) GetMonitorItemsById(c *gin.Context) {
	productId := c.Query("productId")
	osType := c.Query("osType")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", d.dao.SelectMonitorItemsById(productId, osType)))
}
