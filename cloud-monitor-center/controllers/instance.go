package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dbUtils"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceCtl struct{}

func NewInstanceCtl() *InstanceCtl {
	return &InstanceCtl{}
}

func (ctl *InstanceCtl) Page(c *gin.Context) {
	var param forms.InstanceRulePageReqParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.Instance.SelectInstanceRulePage(&param)))
}

func (ctl *InstanceCtl) Unbind(c *gin.Context) {
	var param forms.UnBindRuleParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	err := dbUtils.Tx(&param, service.UnbindInstance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("成功", nil))
}

func (ctl *InstanceCtl) Bind(c *gin.Context) {
	var param forms.InstanceBindRuleDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	err = dbUtils.Tx(&param, service.BindInstance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("成功", nil))
}

func (ctl *InstanceCtl) GetRuleList(c *gin.Context) {
	var param forms.ProductRuleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.Instance.GetRuleListByProductType(&param)))
}
