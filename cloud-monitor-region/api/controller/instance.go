package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	commonUtil "code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceCtl struct {
	dao *dao.InstanceDao
}

func NewInstanceCtl(dao *dao.InstanceDao) *InstanceCtl {
	return &InstanceCtl{dao}
}

// GetPage godoc
// @Summary 产品实例列表
// @Description 产品实例分页
// @Tags InstanceCtl
// @Accept  json
// @Produce  json
// @Param product query string false "产品简称"
// @Param instanceId query string false "实例id"
// @Param instanceName query string false "实例名称"
// @Param statusList query []string false "实例状态"
// @Param extraAttr query string false "其他参数"
// @Param current query  int false "当前页"
// @Param pageSize query int false "页大小"
// @Success 200 {object} vo.PageVO
// @Failure 400 {object} global.Resp
// @Failure 500 {object} global.Resp
// @Router /hawkeye/instance/page [get]
func (ctl *InstanceCtl) GetPage(c *gin.Context) {
	tenantId, _ := commonUtil.GetTenantId(c)
	f := commonService.InstancePageForm{}
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	f.TenantId = tenantId
	instanceService := external.ProductInstanceServiceMap[f.Product]
	page, err := instanceService.GetPage(f, instanceService.(commonService.InstanceStage))
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError("查询失败"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}

// GetInstanceNumByRegion
// @Summary GetInstanceNumByRegion
// @Schemes
// @Description GetInstanceNumByRegion
// @Tags InstanceCtl
// @Accept json
// @Produce json
// @Param product query string false "产品简称"
// @Param instanceId query string false "实例id"
// @Param instanceName query string false "实例名称"
// @Param statusList query []string false "实例状态"
// @Param extraAttr query string false "其他参数"
// @Param current query  int false "当前页"
// @Param pageSize query int false "页大小"
// @Success 200 {object} AlarmInstanceRegionVO
// @Router /hawkeye/instance/getInstanceNum [get]
func (ctl *InstanceCtl) GetInstanceNumByRegion(c *gin.Context) {
	tenantId, _ := commonUtil.GetTenantId(c)
	f := commonService.InstancePageForm{
		TenantId: tenantId,
		Product:  external.ECS,
		Current:  1,
		PageSize: 1000,
	}
	instanceService := external.ProductInstanceServiceMap[external.ECS]
	page, err := instanceService.GetPage(f, instanceService.(commonService.InstanceStage))
	if err != nil {
		logger.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, global.NewError("查询失败"))
		return
	}
	bindNum := ctl.dao.GetInstanceNum(tenantId)
	total := util.If(bindNum > page.Total, bindNum, page.Total)
	vo := &AlarmInstanceRegionVO{Total: total.(int), BindNum: bindNum}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", vo))
}

type AlarmInstanceRegionVO struct {
	Total int `json:"total"`

	/**
	 * 已绑定实例数
	 */
	BindNum int `json:"bindNum"`
}
