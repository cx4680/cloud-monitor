package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	commonUtil "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
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
	f := commonService.InstancePageForm{}
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, f.Product)
	tenantId, iamUserId, err := commonUtil.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	f.TenantId = tenantId
	instanceService := external.ProductInstanceServiceMap[f.Product]
	if instanceService == nil {
		c.JSON(http.StatusBadRequest, global.NewError("该产品未接入"))
		return
	}
	var page *vo.PageVO
	isIamLogin := service.CheckIamLogin(tenantId, iamUserId)
	if isIamLogin {
		f.IamInfo.UserInfo = c.Request.Header.Get("user-info")
		ctl.fillIamInfo(c, &f)
		page, err = instanceService.GetPageByAuth(f, instanceService.(commonService.InstanceStage))
	} else {
		page, err = instanceService.GetPage(f, instanceService.(commonService.InstanceStage))
	}
	if err != nil {
		logger.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, global.NewError("查询失败"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}

func (ctl *InstanceRegionCtl) fillIamInfo(c *gin.Context, f *commonService.InstancePageForm) error {
	header := c.Request.Header
	f.IamInfo.UserInfo = header.Get("user-info")
	sid, _ := c.Cookie("SID")
	f.IamInfo.SID = sid
	f.IamInfo.CurrentTime = header.Get("cs-CurrentTime")
	f.IamInfo.SecureTransport = header.Get("cs-SecureTransport")
	f.IamInfo.SourceIp = header.Get("cs-SourceIp")
	return nil
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
	regionCode := c.Query("region")
	if len(regionCode) == 0 {
		logger.Logger().Error("region不能为空")
		c.JSON(http.StatusOK, global.NewError("region不能为空"))
		return
	}
	tenantId, iamUserId, err := commonUtil.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	f := commonService.InstancePageForm{
		TenantId: tenantId,
		Product:  external.ECS,
		Current:  1,
		PageSize: 1000,
	}
	instanceService := external.ProductInstanceServiceMap[external.ECS]
	var result = &AlarmInstanceRegionVO{}
	var page *vo.PageVO
	isIamLogin := service.CheckIamLogin(tenantId, iamUserId)
	if isIamLogin {
		f.IamInfo.UserInfo = c.Request.Header.Get("user-info")
		ctl.fillIamInfo(c, &f)
		page, err = instanceService.GetPageByAuth(f, instanceService.(commonService.InstanceStage))
	} else {
		page, err = instanceService.GetPage(f, instanceService.(commonService.InstanceStage))
	}
	if err != nil {
		logger.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, global.NewError("查询失败"))
		return
	}
	bindNum := ctl.dao.GetInstanceNum(tenantId, regionCode)
	total := util.If(bindNum > page.Total, bindNum, page.Total)
	result = &AlarmInstanceRegionVO{Total: total.(int), BindNum: bindNum}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", result))
}

type AlarmInstanceRegionVO struct {
	Total int `json:"total"`

	/**
	 * 已绑定实例数
	 */
	BindNum int `json:"bindNum"`
}
