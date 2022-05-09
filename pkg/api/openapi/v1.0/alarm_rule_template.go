package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq/handler"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type AlarmRuleTemplateCtl struct{}

func NewAlarmRuleTemplateCtl() *AlarmRuleTemplateCtl {
	return &AlarmRuleTemplateCtl{}
}

func (ctl *AlarmRuleTemplateCtl) GetProductList(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	list := dao.AlarmRuleTemplate.QueryTemplateProductList(global.DB, tenantId)

	ret := struct {
		RequestId string
		Products  []dao.TemplateProduct
	}{
		RequestId: openapi.GetRequestId(c),
		Products:  list,
	}

	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ret))
}

func (ctl *AlarmRuleTemplateCtl) GetRuleListByProduct(c *gin.Context) {
	productBizId := c.Param("productBizId")
	if len(productBizId) == 0 {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	c.Set(global.ResourceName, productBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}

	ruleList := dao.AlarmRuleTemplate.QueryRuleTemplateListByProduct(global.DB, tenantId, productBizId)
	for i, _ := range ruleList {
		itemList := dao.AlarmItemTemplate.QueryItemListByTemplate(global.DB, ruleList[i].RuleTemplateId)
		cs := make([]string, len(itemList))
		for j, item := range itemList {
			cs[j] = dao.GetExpress2(*item.TriggerCondition)
		}
		ruleList[i].Conditions = cs
	}
	ret := struct {
		RequestId string
		Rules     []dao.AlarmRuleTemplateRule
	}{
		RequestId: openapi.GetRequestId(c),
		Rules:     ruleList,
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ret))
}

func (ctl *AlarmRuleTemplateCtl) Open(c *gin.Context) {
	productBizId := c.Param("productBizId")
	if len(productBizId) == 0 {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	c.Set(global.ResourceName, productBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	userId, err := util.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}

	paramList, errCode := buildRuleReqs(productBizId, tenantId, userId)
	if errCode != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		for _, param := range paramList {
			rel := model.TenantAlarmTemplateRel{
				TenantId:      tenantId,
				TemplateBizId: param.TemplateBizId,
				CreateTime:    util.TimeToFullTimeFmtStr(time.Time{}),
			}
			dao.TenantAlarmTemplateRel.Insert(tx, rel)
			//保存告警规则
			if err = util.Tx(&param, service.CreateRule); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}

	for _, param := range paramList {
		//本Region异步处理
		go handler.HandleMsg(enum.CreateRule, []byte(jsonutil.ToString(param)), false)
	}

	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func (ctl *AlarmRuleTemplateCtl) Close(c *gin.Context) {
	productBizId := c.Param("productBizId")
	if len(productBizId) == 0 {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	c.Set(global.ResourceName, productBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	templateBizIdList := dao.AlarmRuleTemplate.QueryTemplateBizIdListByProductBizId(global.DB, productBizId)
	if len(templateBizIdList) == 0 {
		c.JSON(http.StatusOK, openapi.NewResSuccess(c))
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		dao.TenantAlarmTemplateRel.Delete(tx, tenantId, templateBizIdList)
		for _, templateBizId := range templateBizIdList {
			//delete rules
			ruleIds := dao.AlarmRule.GetRuleIdsByTemplateId(tx, tenantId, templateBizId)
			for _, ruleId := range ruleIds {
				err = dao.AlarmRule.DeleteRule(tx, &form.RuleReqDTO{
					Id:       ruleId,
					TenantId: tenantId,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}

	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func buildRuleReqs(productBizId string, tenantId string, userId string) ([]form.AlarmRuleAddReqDTO, *openapi.ErrorCode) {
	paramList := dao.AlarmRuleTemplate.QueryCreateRuleInfo(global.DB, productBizId)
	if len(paramList) == 0 {
		return nil, openapi.InvalidParameter
	}

	for j, _ := range paramList {
		itemList := dao.AlarmItemTemplate.QueryItemListByTemplate(global.DB, paramList[j].TemplateBizId)
		cs := make([]form.Condition, len(itemList))
		for i, item := range itemList {
			tc := item.TriggerCondition
			cs[i] = form.Condition{
				MetricName:         tc.MetricName,
				MetricCode:         item.MetricCode,
				Period:             tc.Period,
				Times:              tc.Times,
				Statistics:         tc.Statistics,
				ComparisonOperator: tc.ComparisonOperator,
				Threshold:          tc.Threshold,
				Unit:               tc.Unit,
				Labels:             tc.Labels,
				Level:              item.Level,
				SilencesTime:       item.SilencesTime,
				Express:            dao.GetExpress2(*tc),
			}
		}
		paramList[j].Conditions = cs
		paramList[j].GroupList = []string{"-1"}

		p := dao.MonitorProduct.GetMonitorProductByBizId(strconv.Itoa(paramList[j].ProductId))
		resourceList, err := task.GetRemoteProductInstanceList(p.Abbreviation, tenantId)
		if err != nil {
			return nil, openapi.InvalidParameter
		}
		is := make([]*form.InstanceInfo, len(resourceList))
		for i, r := range resourceList {
			is[i] = &form.InstanceInfo{
				InstanceId:   r.InstanceID,
				ZoneCode:     r.ZoneCode,
				RegionCode:   r.RegionCode,
				RegionName:   r.RegionName,
				ZoneName:     r.ZoneName,
				Ip:           r.Ip,
				Status:       "",
				InstanceName: r.InstanceName,
			}
		}
		paramList[j].ResourceList = is

		paramList[j].AlarmHandlerList = []*form.Handler{{
			HandleType: 1,
		}, {
			HandleType: 2,
		}}

		paramList[j].UserId = userId
		paramList[j].TenantId = tenantId
		if errCode := CheckAndFillParam(&paramList[j]); errCode != nil {
			return nil, errCode
		}

		if errCode := addMetricInfo(&paramList[j]); errCode != nil {
			return nil, errCode
		}
	}

	return paramList, nil
}

func addMetricInfo(param *form.AlarmRuleAddReqDTO) *openapi.ErrorCode {
	for i, cond := range param.Conditions {
		metricCode := cond.MetricCode
		item := dao.AlarmRule.GetMonitorItem(metricCode)
		if item.Id == 0 {
			return openapi.InvalidParameter
		}
		param.Conditions[i].Labels = item.Labels
		param.Conditions[i].Unit = item.Unit
		param.Conditions[i].MetricName = item.Name
	}
	return nil
}
