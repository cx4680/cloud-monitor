package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/config"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/global"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/rocketmq/producer"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AlarmRuleCtl struct {
	dao *dao.AlarmRuleDao
}

func NewAlarmRuleCtl(dao *dao.AlarmRuleDao) *AlarmRuleCtl {
	return &AlarmRuleCtl{dao: dao}
}

func (ctl *AlarmRuleCtl) SelectRulePageList(c *gin.Context) {
	var param forms.AlarmPageReqParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, _ := c.Get(global.TenantId)
	param.TenantId = tenantId.(string)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.dao.SelectRulePageList(&param)))
}

func (ctl *AlarmRuleCtl) GetDetail(c *gin.Context) {
	id := c.PostForm("id")
	tenantId, _ := c.Get(global.TenantId)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ctl.dao.GetDetail(id, tenantId.(string))))

}

func (ctl *AlarmRuleCtl) CreateRule(c *gin.Context) {
	var param forms.AlarmRuleAddReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, _ := c.Get(global.TenantId)
	param.TenantId = tenantId.(string)
	userId, _ := c.Get(global.UserId)
	param.UserId = userId.(string)
	addMetricName(&param, ctl)
	id := ctl.dao.SaveRule(&param)
	param.Id = id
	producer.SendMsg(config.GetConfig().Rocketmq.RuleTopic, enums.CreateRule, param)
	c.JSON(http.StatusOK, global.NewSuccess("创建成功", true))
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var param forms.AlarmRuleAddReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, _ := c.Get(global.TenantId)
	param.TenantId = tenantId.(string)
	userId, _ := c.Get(global.UserId)
	param.UserId = userId.(string)
	addMetricName(&param, ctl)
	ctl.dao.UpdateRule(&param)
	producer.SendMsg(config.GetConfig().Rocketmq.RuleTopic, enums.UpdateRule, param)
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func (ctl *AlarmRuleCtl) DeleteRule(c *gin.Context) {
	var param forms.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, _ := c.Get(global.TenantId)
	param.TenantId = tenantId.(string)
	userId, _ := c.Get(global.UserId)
	param.TenantId = userId.(string)
	ctl.dao.DeleteRule(&param)
	producer.SendMsg(config.GetConfig().Rocketmq.RuleTopic, enums.DeleteRule, param)
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func (ctl *AlarmRuleCtl) ChangeRuleStatus(c *gin.Context) {
	var param forms.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	tenantId, _ := c.Get(global.TenantId)
	param.TenantId = tenantId.(string)
	userId, _ := c.Get(global.UserId)
	param.TenantId = userId.(string)
	ctl.dao.UpdateRuleState(&param)
	enum := enums.DisableRule
	if strings.EqualFold(param.Status, dao.ENABLE) {
		enum = enums.EnableRule
	}
	producer.SendMsg(config.GetConfig().Rocketmq.RuleTopic, enum, param)
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func addMetricName(param *forms.AlarmRuleAddReqDTO, ctl *AlarmRuleCtl) {
	item := ctl.dao.GetMonitorItem(param.RuleCondition.MetricName)
	param.RuleCondition.Labels = item.Labels
	param.RuleCondition.Unit = item.Unit
	param.RuleCondition.MonitorItemName = item.Name
}
