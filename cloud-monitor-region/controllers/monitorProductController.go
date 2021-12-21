package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
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
// @Success 200 {object} string
// @Router /hawkeye/monitorProduct/getById [get]
func (mpc *MonitorProductCtl) GetById(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.GetById(id)))
}

func (mpc *MonitorProductCtl) UpdateById(c *gin.Context) {
	var f forms.MonitorProductUpdateForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}

	var p = &models.MonitorProduct{
		Id:         f.Id,
		Name:       f.Name,
		CreateUser: f.CreateUser,
	}
	mpc.dao.UpdateById(p)
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", nil))

}

type Proxy struct {
	*gin.Context
	action ActionInfo
}
type ActionInfo struct {
	Action  string
	Product string
}

func (mpc *MonitorProductCtl) GetById11(c *Proxy) {
	id := c.Query("id")
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", mpc.dao.GetById(id)))
}
