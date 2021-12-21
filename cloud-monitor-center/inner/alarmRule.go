package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dbUtils"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums/sourceType"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
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
	var param = forms.AlarmRuleAddReqDTO{
		SourceType: sourceType.AutoScaling,
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	userId, err := tools.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.UserId = userId
	addMetricName(&param)
	err = dbUtils.Tx(&param, service.CreateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("创建成功", param.Id))
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var param = forms.AlarmRuleAddReqDTO{
		SourceType: sourceType.AutoScaling,
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	userId, err := tools.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.UserId = userId
	addMetricName(&param)
	err = dbUtils.Tx(&param, service.UpdateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func addMetricName(param *forms.AlarmRuleAddReqDTO) {
	item := dao.AlarmRule.GetMonitorItem(param.RuleCondition.MetricName)
	param.RuleCondition.Labels = item.Labels
	param.RuleCondition.Unit = item.Unit
	param.RuleCondition.MonitorItemName = item.Name
}
