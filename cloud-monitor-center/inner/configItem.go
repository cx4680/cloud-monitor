package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConfigItemController struct {
}

func NewConfigItemController() *ConfigItemController {
	return new(ConfigItemController)
}

func (ctl *ConfigItemController) GetItemListById(c *gin.Context) {
	id := c.Query("id")
	if tools.IsBlank(id) {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(id)))
}
