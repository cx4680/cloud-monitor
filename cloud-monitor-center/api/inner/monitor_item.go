package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorItemController struct {
	dao *dao.MonitorItemDao
}

func NewMonitorItemController() *MonitorItemController {
	return &MonitorItemController{dao: dao.MonitorItem}
}

func (ctl *MonitorItemController) GetMonitorItemsById(c *gin.Context) {
	productId := c.Query("productId")
	osType := c.Query("osType")
	if strutil.IsBlank(productId) {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.dao.SelectMonitorItemsById(productId, osType)))
}
