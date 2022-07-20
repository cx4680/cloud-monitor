package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type AlarmRuleCtl struct {
	service *service.AlarmRuleService
}

func NewAlarmRuleCtl() *AlarmRuleCtl {
	return &AlarmRuleCtl{service.NewAlarmRuleService()}
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
	if err = CheckAndFillParam(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	if err = AddMetricInfo(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	err = util.Tx(&param, service.CreateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	go k8s.PrometheusRule.GenerateUserPrometheusRule(tenantId)
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
	if err = CheckAndFillParam(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	err = AddMetricInfo(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	err = util.Tx(&param, service.UpdateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	go k8s.PrometheusRule.GenerateUserPrometheusRule(tenantId)
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
	go k8s.PrometheusRule.GenerateUserPrometheusRule(tenantId)
	c.JSON(http.StatusOK, global.NewSuccess("删除成功", true))
}

func (ctl *AlarmRuleCtl) ChangeRuleStatus(c *gin.Context) {
	var param form.RuleReqDTO
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	if len(param.Status) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("状态值不应为空"))
		return
	}
	c.Set(global.ResourceName, param.Id)
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
	go k8s.PrometheusRule.GenerateUserPrometheusRule(tenantId)
	c.JSON(http.StatusOK, global.NewSuccess("更新成功", true))
}

func CheckAndFillParam(param *form.AlarmRuleAddReqDTO) error {
	if param.Type == constant.AlarmRuleTypeSingleMetric {
		if len(param.MetricCode) == 0 {
			return errors.New("指标不能为空")
		}
		if len(param.SilencesTime) == 0 {
			return errors.New("告警间隔不能为空")
		}
		for i, cond := range param.Conditions {
			param.Conditions[i].MetricCode = param.MetricCode
			param.Conditions[i].SilencesTime = param.SilencesTime
			if cond.Level == 0 {
				return errors.New("告警级别不能为空")
			}

			if cond.Period == 0 {
				return errors.New("数据周期不能为空")
			}
			if cond.Times == 0 {
				return errors.New("持续周期不能为空")
			}
		}
	} else {
		if param.Level == 0 {
			return errors.New("告警级别不能为空")
		}
		if len(param.SilencesTime) == 0 {
			return errors.New("告警间隔不能为空")
		}
		if param.Period == 0 {
			return errors.New("数据周期不能为空")
		}
		if param.Times == 0 {
			return errors.New("持续周期不能为空")
		}
		if param.Combination != 1 && param.Combination != 2 {
			return errors.New("告警条件关系有误")
		}
	}
	return nil
}

func AddMetricInfo(param *form.AlarmRuleAddReqDTO) error {
	for i, cond := range param.Conditions {
		metricCode := cond.MetricCode
		item := dao.AlarmRule.GetMonitorItem(metricCode)
		if item.Id == 0 {
			return errors.New("指标不存在")
		}
		param.Conditions[i].Labels = item.Labels
		param.Conditions[i].Unit = item.Unit
		param.Conditions[i].MetricName = item.Name
	}
	return nil
}
