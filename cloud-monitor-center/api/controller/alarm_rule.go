package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type AlarmRuleCtl struct {
}

func NewAlarmRuleCtl() *AlarmRuleCtl {
	return &AlarmRuleCtl{}
}

func (ctl *AlarmRuleCtl) SelectRulePageList(c *gin.Context) {
	var param form.AlarmPageReqParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
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
	c.Set(global.ResourceName, id)
	tenantId, _ := util.GetTenantId(c)
	detail, _ := dao.AlarmRule.GetDetail(global.DB, id, tenantId)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", detail))

}

func (ctl *AlarmRuleCtl) CreateRule(c *gin.Context) {
	var param = form.AlarmRuleAddReqDTO{
		SourceType: source_type.Front,
	}
	CreateRule(c, param)
}

func CreateRule(c *gin.Context, param form.AlarmRuleAddReqDTO) {
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	userId, err := util.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.UserId = userId
	err = addMetricName(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	err = util.Tx(&param, service.CreateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Set(global.ResourceName, param.Id)
	c.JSON(http.StatusOK, global.NewSuccess("创建成功", param.Id))
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var param form.AlarmRuleAddReqDTO
	UpdateRule(c, param)
}

func UpdateRule(c *gin.Context, param form.AlarmRuleAddReqDTO) {
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	c.Set(global.ResourceName, param.Id)
	param.TenantId = tenantId
	userId, err := util.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.UserId = userId
	err = addMetricName(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	err = util.Tx(&param, service.UpdateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func (ctl *AlarmRuleCtl) DeleteRule(c *gin.Context) {
	var param form.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.Id)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	err = util.Tx(&param, service.DeleteRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
}

func (ctl *AlarmRuleCtl) ChangeRuleStatus(c *gin.Context) {
	var param form.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	c.Set(global.ResourceName, param.Id)
	if len(param.Status) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("状态值不应为空"))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	err = util.Tx(&param, service.ChangeRuleStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func addMetricName(param *form.AlarmRuleAddReqDTO) error {
	item := dao.AlarmRule.GetMonitorItem(param.RuleCondition.MetricName)
	if item == nil {
		return errors.New("指标不存在")
	}
	param.RuleCondition.Labels = item.Labels
	param.RuleCondition.Unit = item.Unit
	param.RuleCondition.MonitorItemName = item.Name
	return nil
}
