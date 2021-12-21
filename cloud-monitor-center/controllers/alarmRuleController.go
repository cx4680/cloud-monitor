package controllers

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

func (ctl *AlarmRuleCtl) SelectRulePageList(c *gin.Context) {
	var param forms.AlarmPageReqParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.AlarmRule.SelectRulePageList(&param)))
}

func (ctl *AlarmRuleCtl) GetDetail(c *gin.Context) {
	id := c.PostForm("id")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("缺少id"))
		return
	}
	tenantId, _ := tools.GetTenantId(c)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", dao.AlarmRule.GetDetail(global.DB, id, tenantId)))

}

func (ctl *AlarmRuleCtl) CreateRule(c *gin.Context) {
	var param = forms.AlarmRuleAddReqDTO{
		SourceType: sourceType.Front,
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
	var param forms.AlarmRuleAddReqDTO
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

func (ctl *AlarmRuleCtl) DeleteRule(c *gin.Context) {
	var param forms.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	err = dbUtils.Tx(&param, service.DeleteRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func (ctl *AlarmRuleCtl) ChangeRuleStatus(c *gin.Context) {
	var param forms.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := tools.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
	}
	param.TenantId = tenantId
	err = dbUtils.Tx(&param, service.ChangeRuleStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
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
