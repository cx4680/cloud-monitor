package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/message_center"
	external "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/region"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
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
	var noticeChannelList = message_center.NewService().GetRemoteChannels()
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", noticeChannelList))
}

func (ctl *ConfigItemCtl) CheckCloudLogin(c *gin.Context) {
	tenantId, iamUserId, err := util.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	isOpen, err := service.CheckIamDirectory(tenantId)
	if err != nil {
		logger.Logger().Errorf("IamLogin接口错误：%v", err)
	}
	if !isOpen || strutil.IsBlank(iamUserId) || iamUserId == tenantId {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", true))
	} else {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", false))
	}
}
