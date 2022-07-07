package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/region"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq/handler"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type AlarmRuleTemplateCtl struct {
	MessageCenterSvc *message_center.Service
	RegionSvc        *region.ExternService
}

func NewAlarmRuleTemplateCtl(MessageCenterSvc *message_center.Service, RegionSvc *region.ExternService) *AlarmRuleTemplateCtl {
	return &AlarmRuleTemplateCtl{MessageCenterSvc: MessageCenterSvc, RegionSvc: RegionSvc}
}

func (ctl *AlarmRuleTemplateCtl) GetProductList(c *gin.Context) {
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	list := dao.AlarmRuleTemplate.QueryTemplateProductList(global.DB, tenantId)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", list))
}

func (ctl *AlarmRuleTemplateCtl) GetRuleListByProduct(c *gin.Context) {
	productBizId := c.Query("productBizId")
	if len(productBizId) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("参数错误"))
		return
	}
	c.Set(global.ResourceName, productBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}

	ruleList := dao.AlarmRuleTemplate.QueryRuleTemplateListByProduct(global.DB, tenantId, productBizId)
	for i, v := range ruleList {
		//判断告警模板是否已创建规则，若已创建则使用规则ID查询规则名称触发条件
		if v.Type == 1 {
			ruleDetail, _ := dao.AlarmRule.GetDetail(global.DB, v.RuleId, tenantId)
			cs := make([]string, len(ruleDetail.RuleConditions))
			for j, condition := range ruleDetail.RuleConditions {
				cs[j] = dao.GetExpress(*condition)
			}
			ruleList[i].RuleName = ruleDetail.RuleName
			ruleList[i].Conditions = cs
		} else if v.Type == 0 {
			itemList := dao.AlarmItemTemplate.QueryItemListByTemplate(global.DB, v.RuleTemplateId)
			cs := make([]string, len(itemList))
			for j, item := range itemList {
				cs[j] = dao.GetExpress2(*item.TriggerCondition)
			}
			ruleList[i].Conditions = cs
		}
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", ruleList))
}

func (ctl *AlarmRuleTemplateCtl) Open(c *gin.Context) {
	productBizId := c.PostForm("productBizId")
	if len(productBizId) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.Set(global.ResourceName, productBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	userId, err := util.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}

	paramList, err := ctl.buildRuleReqs(productBizId, tenantId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
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

	c.JSON(http.StatusOK, global.NewSuccess("开启成功", nil))
}

func (ctl *AlarmRuleTemplateCtl) Close(c *gin.Context) {
	productBizId := c.PostForm("productBizId")
	if len(productBizId) == 0 {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.Set(global.ResourceName, productBizId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(err.Error()))
		return
	}
	templateBizIdList := dao.AlarmRuleTemplate.QueryTemplateBizIdListByProductBizId(global.DB, productBizId)
	if len(templateBizIdList) == 0 {
		c.JSON(http.StatusOK, global.NewSuccess("关闭成功", nil))
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
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, global.NewSuccess("关闭成功", nil))
}

func (ctl *AlarmRuleTemplateCtl) buildRuleReqs(productBizId string, tenantId string, userId string) ([]form.AlarmRuleAddReqDTO, error) {
	paramList := dao.AlarmRuleTemplate.QueryCreateRuleInfo(global.DB, productBizId)
	if len(paramList) == 0 {
		return nil, errors.NewBusinessError("该产品下无规则")
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

		list, err := ctl.RegionSvc.GetRegionList(tenantId)
		if err != nil {
			return nil, err
		}
		var resourceList []*model.AlarmInstance
		for _, regionObj := range list {
			regionResources, err := ctl.getResourceByRegion(regionObj.Code, p.Abbreviation, tenantId)
			if err != nil {
				continue
			}
			if len(regionResources) > 0 {
				resourceList = append(resourceList, regionResources...)
			}
		}

		if err != nil {
			return nil, errors.NewBusinessError("实例获取失败")
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

		//动态读取通知方式
		var noticeChannelList = ctl.MessageCenterSvc.GetRemoteChannels()
		if len(noticeChannelList) > 0 {
			handlers := make([]*form.Handler, len(noticeChannelList))
			for i, channel := range noticeChannelList {
				handlers[i] = &form.Handler{HandleType: int(channel.Data)}
			}
			paramList[j].AlarmHandlerList = handlers
		}

		paramList[j].UserId = userId
		paramList[j].TenantId = tenantId
		if err = CheckAndFillParam(&paramList[j]); err != nil {
			return nil, err
		}
		if err = AddMetricInfo(&paramList[j]); err != nil {
			return nil, err
		}
	}

	return paramList, nil
}

func (ctl *AlarmRuleTemplateCtl) getResourceByRegion(regionId, abbreviation, tenantId string) ([]*model.AlarmInstance, error) {
	respStr, err := httputil.HttpGet("http://cloud-monitor." + regionId + ".intranet.cecloudcs.com/inner/monitorResource/list?abbreviation" + abbreviation + "&tenantId=" + tenantId)
	if err != nil {
		return nil, err
	}
	var r global.Resp
	err = jsonutil.ToObjectWithError(respStr, &r)
	if err != nil {
		return nil, err
	}
	if r.Module == nil {
		logger.Logger().Infof("get remote resource empty. regionId=%s, abbreviation=%s, tenantId=%s", regionId, abbreviation, tenantId)
		return nil, nil
	}
	list, ok := r.Module.([]interface{})
	if ok {
		logger.Logger().Infof("get remote resource fail, return data type error. regionId=%s, abbreviation=%s, tenantId=%s", regionId, abbreviation, tenantId)
		return nil, nil
	}
	rets := make([]*model.AlarmInstance, len(list))
	for i, temp := range list {
		var ai model.AlarmInstance
		err = jsonutil.ToObjectWithError(jsonutil.ToString(temp), &ai)
		if err != nil {
			logger.Logger().Errorf("get remote resource fail. regionId=%s, abbreviation=%s, tenantId=%s, err=%v", regionId, abbreviation, tenantId, err)
			continue
		}
		rets[i] = &ai
	}
	return rets, nil
}
