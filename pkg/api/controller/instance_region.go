package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	commonUtil "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceRegionCtl struct {
	dao *dao.InstanceDao
}

func NewInstanceRegionCtl(dao *dao.InstanceDao) *InstanceRegionCtl {
	return &InstanceRegionCtl{dao}
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
func (ctl *InstanceRegionCtl) GetPage(c *gin.Context) {
	tenantId, _ := commonUtil.GetTenantId(c)
	f := commonService.InstancePageForm{}
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, f.Product)
	f.TenantId = tenantId
	instanceService := external.ProductInstanceServiceMap[f.Product]
	if instanceService == nil {
		c.JSON(http.StatusBadRequest, global.NewError("该产品未接入"))
		return
	}
	page, err := instanceService.GetPage(f, instanceService.(commonService.InstanceStage))
	if err != nil {
		logger.Logger().Error(err)
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
func (ctl *InstanceRegionCtl) GetInstanceNumByRegion(c *gin.Context) {
	tenantId, _ := commonUtil.GetTenantId(c)
	regionCode := c.Query("region")
	if len(regionCode) == 0 {
		logger.Logger().Error("region不能为空")
		c.JSON(http.StatusOK, global.NewError("region不能为空"))
		return
	}
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
	bindNum := ctl.dao.GetInstanceNum(tenantId, regionCode)
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
