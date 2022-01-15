package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorProductCtl struct {
	dao *dao.MonitorProductDao
}

func NewMonitorProductCtl(dao *dao.MonitorProductDao) *MonitorProductCtl {
	return &MonitorProductCtl{dao}
}

func (mpc *MonitorProductCtl) GetAllMonitorProducts(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.SelectMonitorProductList()))
}
