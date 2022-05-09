package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ConfigItemController struct {
}

func NewConfigItemController() *ConfigItemController {
	return new(ConfigItemController)
}

func (ctl *ConfigItemController) GetItemListById(c *gin.Context) {
	idstr := c.Query("id")
	if strutil.IsBlank(idstr) {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	id, err := strconv.ParseInt(idstr, 0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.Set(global.ResourceName, id)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(id)))
}
