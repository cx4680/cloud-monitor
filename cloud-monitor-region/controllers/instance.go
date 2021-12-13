package controllers

import (
	commonGlobal "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
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
// @Params product query string false "产品简称"
// @Params instanceId query string false "实例id"
// @Params instanceName query string false "实例名称"
// @Params statusList query []string false "实例状态"
// @Params extraAttr query map[string]string false "其他参数"
// @Params current query  int false "当前页"
// @Params pageSize query int false "页大小"
// @Success 200 {object} vo.PageVO
// @Failure 400 {object} global.Resp
// @Failure 500 {object} global.Resp
// @Router /hawkeye/instance/page [get]
func (ctl *InstanceCtl) GetPage(c *gin.Context) {
	tenantId, _ := c.Get(commonGlobal.TenantId)
	form := service.InstancePageForm{}
	if err := c.ShouldBindQuery(form); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	form.TenantId = tools.ToString(tenantId)
	instanceService := external.ProductInstanceServiceMap[form.Product]
	page, err := instanceService.GetPage(form, instanceService.(service.InstanceStage))
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
// @Params product query string false "产品简称"
// @Params instanceId query string false "实例id"
// @Params instanceName query string false "实例名称"
// @Params statusList query []string false "实例状态"
// @Params extraAttr query map[string]string false "其他参数"
// @Params current query  int false "当前页"
// @Params pageSize query int false "页大小"
// @Success 200 {object} AlarmInstanceRegionVO
// @Router /hawkeye/instance/page [get]
func (ctl *InstanceCtl) GetInstanceNumByRegion(c *gin.Context) {
	tenantId, _ := c.Get(commonGlobal.TenantId)
	form := service.InstancePageForm{}
	if err := c.ShouldBindQuery(form); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	form.TenantId = tools.ToString(tenantId)
	instanceService := external.ProductInstanceServiceMap[external.ECS]
	page, err := instanceService.GetPage(form, instanceService.(service.InstanceStage))
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError("查询失败"))
		return
	}
	bindNum := ctl.dao.GetInstanceNum(tenantId.(string))
	total := utils.If(bindNum > page.Total, bindNum, page.Total)
	vo := &AlarmInstanceRegionVO{Total: total.(int), BindNum: bindNum}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", vo))
}

type AlarmInstanceRegionVO struct {
	Total int

	/**
	 * 已绑定实例数
	 */
	BindNum int
}
