package controllers

import (
	"business-common/dao"
	"business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/config"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/global"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/rocketmq/producer"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/validator/translate"
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
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	ctl.dao.UnbindInstance(&param)
	producer.SendMsg(config.GetConfig().Rocketmq.RuleTopic, enums.UnbindRule, param)
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
	producer.SendMsg(config.GetConfig().Rocketmq.RuleTopic, enums.BindRule, param)
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
