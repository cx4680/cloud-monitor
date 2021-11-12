package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/global"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceCtl struct {
	dao *dao.InstanceDao
}

func NewInstanceCtl(dao *dao.InstanceDao) *InstanceCtl {
	return &InstanceCtl{dao: dao}
}

func (ctl *InstanceCtl) Page(c *gin.Context) {
	var param forms.InstanceRulePageReqParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.dao.SelectInstanceRulePage(&param)))
}

func (ctl *InstanceCtl) Unbind(c *gin.Context) {
	var param forms.UnBindRuleParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	ctl.dao.UnbindInstance(&param)
	mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enums.UnbindRule, param)
	c.JSON(http.StatusOK, global.NewSuccess("成功", nil))
}

func (ctl *InstanceCtl) Bind(c *gin.Context) {
	var param forms.InstanceBindRuleDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, _ := c.Get(global.TenantId)
	param.TenantId = tenantId.(string)
	ctl.dao.BindInstance(&param)
	mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enums.BindRule, param)
	c.JSON(http.StatusOK, global.NewSuccess("成功", nil))
}

func (ctl *InstanceCtl) GetRuleList(c *gin.Context) {
	var param forms.ProductRuleParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, _ := c.Get(global.TenantId)
	param.TenantId = tenantId.(string)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.dao.GetRuleListByProductType(&param)))
}
