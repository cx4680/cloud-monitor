package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceCtl struct{}

func NewInstanceCtl() *InstanceCtl {
	return &InstanceCtl{}
}

func (ctl *InstanceCtl) Page(c *gin.Context) {
	var param form.InstanceRulePageReqParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.InstanceId)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.Instance.SelectInstanceRulePage(&param)))
}

func (ctl *InstanceCtl) Unbind(c *gin.Context) {
	var param form.UnBindRuleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.InstanceId)
	tenantId, err2 := util.GetTenantId(c)
	if err2 != nil {
		logger.Logger().Info("tenantId is nil")
		c.JSON(http.StatusBadRequest, global.NewError(err2.Error()))
		return
	}
	param.TenantId = tenantId
	err := util.Tx(&param, service.UnbindInstance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("成功", nil))
}

func (ctl *InstanceCtl) Bind(c *gin.Context) {
	var param form.InstanceBindRuleDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.InstanceId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	err = util.Tx(&param, service.BindInstance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("成功", nil))
}

func (ctl *InstanceCtl) GetRuleList(c *gin.Context) {
	var param form.ProductRuleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.InstanceId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.Instance.GetRuleListByProductType(&param)))
}
