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

// GetById
// @Summary GetById
// @Schemes
// @Description GetById
// @Tags MonitorProductCtl
// @Accept json
// @Produce json
// @Param id query  string true "id"
// @Success 200 {object} model.MonitorProduct
// @Router /hawkeye/monitorProduct/getById [get]
func (mpc *MonitorProductCtl) GetById(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.GetById(id)))
}

func (mpc *MonitorProductCtl) GetAllMonitorProducts(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.SelectMonitorProductList()))
}

type Proxy struct {
	*gin.Context
	action ActionInfo
}
type ActionInfo struct {
	Action  string
	Product string
}
