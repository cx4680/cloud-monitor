package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	external "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/region"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConfigItemCtl struct {
	service *external.ExternService
}

func NewConfigItemCtl() *ConfigItemCtl {
	return &ConfigItemCtl{external.NewExternService()}
}

func (ctl *ConfigItemCtl) GetStatisticalPeriodList(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(dao.StatisticalPeriodPid)))
}

func (ctl *ConfigItemCtl) GetContinuousCycleList(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(dao.ContinuousCyclePid)))
}
func (ctl *ConfigItemCtl) GetStatisticalMethodsList(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(dao.StatisticalMethodsPid)))
}
func (ctl *ConfigItemCtl) GetComparisonMethodList(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(dao.ComparisonMethodPid)))
}
func (ctl *ConfigItemCtl) GetOverviewItemList(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(dao.OverviewItemPid)))
}
func (ctl *ConfigItemCtl) GetRegionList(c *gin.Context) {
	tenantId, _ := c.Get(global.TenantId)
	list, err := ctl.service.GetRegionList(tenantId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", list))
}
func (ctl *ConfigItemCtl) GetMonitorRange(c *gin.Context) {
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.ConfigItem.GetConfigItemList(dao.MonitorRange)))
}

func (ctl *ConfigItemCtl) GetNoticeChannel(c *gin.Context) {
	if config.Cfg.Common.MsgIsOpen == config.MsgClose {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", nil))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", global.NoticeChannelList))
}
