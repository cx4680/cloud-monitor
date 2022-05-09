package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
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
	var param form.MonitorItemParam
	err := c.ShouldBindQuery(&param)
	if err != nil || strutil.IsBlank(param.ProductBizId) {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.Set(global.ResourceName, param.ProductBizId)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.dao.GetMonitorItem(param.ProductBizId, param.OsType, param.Display)))
}
