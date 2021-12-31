package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlarmRuleCtl struct {
}

func NewAlarmRuleCtl() *AlarmRuleCtl {
	return &AlarmRuleCtl{}
}

func (ctl *AlarmRuleCtl) CreateRule(c *gin.Context) {
	var param = form.AlarmRuleAddReqDTO{
		SourceType: source_type.AutoScaling,
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	userId, err := util.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.UserId = userId
	addMetricName(&param)
	err = util.Tx(&param, service.CreateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("创建成功", param.Id))
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var param = form.AlarmRuleAddReqDTO{
		SourceType: source_type.AutoScaling,
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	userId, err := util.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.UserId = userId
	addMetricName(&param)
	err = util.Tx(&param, service.UpdateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func addMetricName(param *form.AlarmRuleAddReqDTO) {
	item := dao.AlarmRule.GetMonitorItem(param.RuleCondition.MetricName)
	param.RuleCondition.Labels = item.Labels
	param.RuleCondition.Unit = item.Unit
	param.RuleCondition.MonitorItemName = item.Name
}
