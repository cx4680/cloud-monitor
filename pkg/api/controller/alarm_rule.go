package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	global2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	util2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
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
		c.JSON(http.StatusBadRequest, global2.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	c.JSON(http.StatusOK, global2.NewSuccess("查询成功", dao.AlarmRule.SelectRulePageList(&param)))
}

func (ctl *AlarmRuleCtl) GetDetail(c *gin.Context) {
	id := c.PostForm("id")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, global2.NewError("缺少id"))
		return
	}
	tenantId, _ := util2.GetTenantId(c)
	detail, _ := dao.AlarmRule.GetDetail(global2.DB, id, tenantId)
	c.JSON(http.StatusOK, global2.NewSuccess("查询成功", detail))

}

func (ctl *AlarmRuleCtl) CreateRule(c *gin.Context) {
	var param = form.AlarmRuleAddReqDTO{
		SourceType: source_type.Front,
	}
	CreateRule(c, param)
}

func CreateRule(c *gin.Context, param form.AlarmRuleAddReqDTO) {
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	userId, err := util2.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	param.UserId = userId
	err = addMetricName(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	err = util2.Tx(&param, service.CreateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, global2.NewSuccess("创建成功", param.Id))
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var param form.AlarmRuleAddReqDTO
	UpdateRule(c, param)
}

func UpdateRule(c *gin.Context, param form.AlarmRuleAddReqDTO) {
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	userId, err := util2.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	param.UserId = userId
	err = addMetricName(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	err = util2.Tx(&param, service.UpdateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, global2.NewSuccess("更新成功", true))
}

func (ctl *AlarmRuleCtl) DeleteRule(c *gin.Context) {
	var param form.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	err = util2.Tx(&param, service.DeleteRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global2.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global2.NewSuccess("删除成功", true))
}

func (ctl *AlarmRuleCtl) ChangeRuleStatus(c *gin.Context) {
	var param form.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(translate.GetErrorMsg(err)))
		return
	}
	if len(param.Status) == 0 {
		c.JSON(http.StatusBadRequest, global2.NewError("状态值不应为空"))
		return
	}
	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global2.NewError(err.Error()))
		return
	}
	param.TenantId = tenantId
	err = util2.Tx(&param, service.ChangeRuleStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global2.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global2.NewSuccess("更新成功", true))
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
