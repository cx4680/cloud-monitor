package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/task"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorResourceController struct {
}

func NewMonitorResourceController() *MonitorResourceController {
	return new(MonitorResourceController)
}

func (ctl *MonitorResourceController) GetProductInstanceList(c *gin.Context) {
	abb := c.Query("abbreviation")
	tenantId := c.Query("tenantId")
	list, err := task.GetRemoteProductInstanceList(abb, tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", list))
}
