package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	commonError "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq/handler"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
)

type AlarmRuleCtl struct {
}

func NewAlarmRuleCtl() *AlarmRuleCtl {
	return &AlarmRuleCtl{}
}

func (ctl *AlarmRuleCtl) SelectRulePageList(c *gin.Context) {
	reqParam := AlarmPageReqParam{
		PageNumber: 1,
		PageSize:   10,
	}
	if err := c.ShouldBindQuery(&reqParam); err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	pageParam := form.AlarmPageReqParam{
		RuleName: reqParam.RuleName,
		Status:   reqParam.Status,
		TenantId: tenantId,
		PageSize: reqParam.PageSize,
		Current:  reqParam.PageNumber,
	}
	pageVo := dao.AlarmRule.SelectRulePageList(&pageParam)
	var ruleList []RuleInfo
	if pageVo.Records != nil {
		listVo := pageVo.Records.([]form.AlarmRulePageDTO)
		for _, ruleVo := range listVo {
			product := dao.MonitorProduct.GetByName(global.DB, ruleVo.ProductType)
			ruleInfo := RuleInfo{
				Name:                ruleVo.Name,
				MonitorType:         ruleVo.MonitorType,
				ProductType:         ruleVo.ProductType,
				ResourceNum:         ruleVo.InstanceNum,
				Status:              ruleVo.Status,
				RuleId:              ruleVo.RuleId,
				ProductAbbreviation: product.Abbreviation,
			}
			ruleList = append(ruleList, ruleInfo)
		}
	}
	page := AlarmRulePageDTO{
		ResCommonPage: *openapi.NewResCommonPage(c, pageVo),
		Rules:         ruleList,
	}
	c.JSON(http.StatusOK, page)
}

func (ctl *AlarmRuleCtl) GetDetail(c *gin.Context) {
	ruleId := c.Param("RuleId")
	c.Set(global.ResourceName, ruleId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		logger.Logger().Info("tenantId is nil")
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	detailVo, err := dao.AlarmRule.GetDetail(global.DB, ruleId, tenantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
		return
	}
	res := RuleDetail{
		RequestId: openapi.GetRequestId(c),
	}
	if detailVo.Scope == "INSTANCE" {
		res.RuleInfo.Scope = "Resource"
	} else {
		res.RuleInfo.Scope = detailVo.Scope
	}
	res.RuleInfo.RuleId = detailVo.Id
	res.RuleInfo.Type = detailVo.Type
	res.RuleInfo.Status = detailVo.Status
	res.RuleInfo.MonitorType = detailVo.MonitorType
	res.RuleInfo.ProductType = detailVo.ProductType
	res.RuleInfo.RuleName = detailVo.RuleName
	res.RuleInfo.SilencesTime = detailVo.SilencesTime
	res.RuleInfo.AlarmLevel = detailVo.AlarmLevel
	res.RuleInfo.Combination = detailVo.Combination
	res.RuleInfo.Describe = detailVo.Describe
	res.RuleInfo.Period = detailVo.Period
	res.RuleInfo.Times = detailVo.Times

	for _, c := range detailVo.RuleConditions {
		res.RuleInfo.RuleConditions = append(res.RuleInfo.RuleConditions, struct {
			MetricCode         string
			MetricName         string
			Period             int
			Times              int
			Statistics         string
			ComparisonOperator string
			Threshold          string
			Unit               string
			Labels             string
			Level              uint8
			SilencesTime       string
			Express            string
		}{MetricCode: c.MetricCode, MetricName: c.MetricName, Period: c.Period,
			Times: c.Times, Statistics: c.Statistics, ComparisonOperator: c.ComparisonOperator, Threshold: fmt.Sprintf("%.2f", c.Threshold), Unit: c.Unit, Labels: c.Labels, Level: c.Level, SilencesTime: c.SilencesTime, Express: c.Express})
	}

	for _, groupVo := range detailVo.NoticeGroups {
		res.RuleInfo.NoticeGroups = append(res.RuleInfo.NoticeGroups, struct {
			Id   string
			Name string
		}{Id: groupVo.Id, Name: groupVo.Name})
	}
	for _, instanceVo := range detailVo.InstanceList {
		res.RuleInfo.ResourceIds = append(res.RuleInfo.ResourceIds, struct {
			ResourceId string
		}{ResourceId: instanceVo.InstanceId})
	}
	for _, alarmHandlerVo := range detailVo.AlarmHandlerList {
		res.RuleInfo.NoticeChannels = append(res.RuleInfo.NoticeChannels, struct {
			HandlerType int
		}{HandlerType: alarmHandlerVo.HandleType})
	}
	c.JSON(http.StatusOK, res)
}

func (ctl *AlarmRuleCtl) CreateRule(c *gin.Context) {
	var reqParam AlarmRuleCreateReqDTO
	addForm, errCode := buildAlarmRuleReqParam(c, &reqParam, nil)
	if errCode != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(errCode, c))
		return
	}
	errCode = CheckAndFillParam(addForm)
	if errCode != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(errCode, c))
		return
	}
	err := util.Tx(addForm, service.CreateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	//本Region异步处理
	go handler.HandleMsg(enum.CreateRule, []byte(jsonutil.ToString(addForm)), false)
	res := struct {
		RequestId string
		RuleId    string
	}{
		RequestId: openapi.GetRequestId(c),
		RuleId:    addForm.Id,
	}
	c.Set(global.ResourceName, addForm.Id)
	c.JSON(http.StatusOK, res)
}

func buildAlarmRuleReqParam(c *gin.Context, createParam *AlarmRuleCreateReqDTO, updateParam *AlarmRuleUpdateReqDTO) (*form.AlarmRuleAddReqDTO, *openapi.ErrorCode) {
	var param = &AlarmRuleCreateReqDTO{}
	var productInfo = &model.MonitorProduct{}
	if createParam != nil {
		if err := c.ShouldBindJSON(&createParam); err != nil {
			logger.Logger().Info(err)
			return nil, openapi.GetErrorCode(err)
		}
		param = createParam
		productInfo = dao.MonitorProduct.GetByAbbreviation(global.DB, param.ProductAbbreviation)
		if productInfo == nil || len(productInfo.BizId) == 0 {
			return nil, openapi.ProductAbbreviationInvalid
		}
	} else if updateParam != nil {
		if err := c.ShouldBindJSON(updateParam); err != nil {
			logger.Logger().Info(err)
			return nil, openapi.GetErrorCode(err)
		}
		param = &AlarmRuleCreateReqDTO{
			Type:              updateParam.Type,
			Scope:             updateParam.Scope,
			Resources:         updateParam.Resources,
			RuleName:          updateParam.RuleName,
			NoticeChannels:    updateParam.NoticeChannels,
			TriggerConditions: updateParam.TriggerConditions,
			SilencesTime:      updateParam.SilencesTime,
			AlarmLevel:        updateParam.AlarmLevel,
			GroupList:         updateParam.GroupList,
			MetricCode:        updateParam.MetricCode,
			Period:            updateParam.Period,
			Times:             updateParam.Times,
			Combination:       updateParam.Combination,
		}
	}
	if len(param.Resources) == 0 {
		return nil, openapi.MissingResources
	}
	nameMatched, err := regexp.MatchString("^[a-z][a-z0-9_]{0,14}[a-z0-9]$", param.RuleName)
	if !nameMatched {
		return nil, openapi.RuleNameInvalid
	}

	tenantId, err := util.GetTenantId(c)
	if err != nil {
		logger.Logger().Info("tenantId is nil")
		return nil, openapi.MissingParameter
	}
	userId, err := util.GetUserId(c)
	if err != nil {
		logger.Logger().Info("userId is nil")
		return nil, openapi.MissingParameter
	}

	productBizId, _ := strconv.Atoi(productInfo.BizId)
	var scope string
	if param.Scope == "Resource" {
		scope = "INSTANCE"
	} else {
		scope = param.Scope
	}

	addForm := &form.AlarmRuleAddReqDTO{
		Type:         param.Type,
		RuleName:     param.RuleName,
		MonitorType:  "云产品监控",
		ProductType:  productInfo.Name,
		ProductId:    productBizId,
		Scope:        scope,
		TenantId:     tenantId,
		UserId:       userId,
		SilencesTime: "3h",
		Level:        param.AlarmLevel,
		MetricCode:   param.MetricCode,
		SourceType:   source_type.Front,
	}
	for _, resource := range param.Resources {
		addForm.ResourceList = append(addForm.ResourceList, &form.InstanceInfo{
			InstanceId: resource.ResourceId,
		})
	}
	for _, channel := range param.NoticeChannels {
		addForm.AlarmHandlerList = append(addForm.AlarmHandlerList, &form.Handler{HandleType: channel.HandlerType})
	}
	addForm.GroupList = param.GroupList

	if param.Type == constant.AlarmRuleTypeSingleMetric {
		err := buildSingleAlarmRuleReqParam(param, addForm)
		if err != nil {
			return nil, err
		}
	}
	if param.Type == constant.AlarmRuleTypeMultipleMetric {
		err := buildMultipleAlarmRuleReqParam(param, addForm)
		if err != nil {
			return nil, err
		}
	}
	return addForm, nil
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var reqParam AlarmRuleUpdateReqDTO
	updateForm, errCode := buildAlarmRuleReqParam(c, nil, &reqParam)
	if errCode != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(errCode, c))
		return
	}
	ruleId := c.Param("RuleId")
	c.Set(global.ResourceName, ruleId)
	updateForm.Id = ruleId
	errCode = CheckAndFillParam(updateForm)
	if errCode != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(errCode, c))
		return
	}
	err := util.Tx(updateForm, service.UpdateRule)
	if err != nil {
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
			return
		}
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	go handler.HandleMsg(enum.UpdateRule, []byte(jsonutil.ToString(updateForm)), false)
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func (ctl *AlarmRuleCtl) DeleteRule(c *gin.Context) {
	ruleId := c.Param("RuleId")
	c.Set(global.ResourceName, ruleId)
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	reqParam := form.RuleReqDTO{
		Id:       ruleId,
		TenantId: tenantId,
	}
	err = util.Tx(&reqParam, service.DeleteRule)
	if err != nil {
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
			return
		}
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	go handler.HandleMsg(enum.DeleteRule, []byte(jsonutil.ToString(&reqParam)), false)
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func (ctl *AlarmRuleCtl) ChangeRuleStatus(c *gin.Context) {
	ruleId := c.Param("RuleId")
	c.Set(global.ResourceName, ruleId)
	var status StatusBody
	if err := c.ShouldBindJSON(&status); err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		logger.Logger().Info("tenantId is nil")
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	reqParam := form.RuleReqDTO{
		TenantId: tenantId,
		Id:       ruleId,
		Status:   status.Enable,
	}
	err = util.Tx(&reqParam, service.ChangeRuleStatus)
	if err != nil {
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
			return
		}
		logger.Logger().Info(err)
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	go handler.HandleMsg(enum.ChangeStatus, []byte(jsonutil.ToString(&reqParam)), false)
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func checkMetricName(metricCode string) (*model.MonitorItem, error) {
	item := dao.AlarmRule.GetMonitorItem(metricCode)
	if item == nil || len(item.MetricName) == 0 {
		return nil, errors.New("指标不存在")
	}
	return item, nil
}

func buildSingleAlarmRuleReqParam(param *AlarmRuleCreateReqDTO, addForm *form.AlarmRuleAddReqDTO) *openapi.ErrorCode {
	monitorItem, err := checkMetricName(param.MetricCode)
	if err != nil {
		return openapi.MetricCodeInvalid
	}
	for _, cond := range param.TriggerConditions {
		matched, err := regexp.MatchString("^[0-9\\.]+$", cond.Threshold)
		if !matched {
			return openapi.ThresholdInvalid
		}
		threshold, err := strconv.ParseFloat(cond.Threshold, 64)
		if err != nil {
			logger.Logger().Infof("Threshold is parsefloat error:%v", err)
			return openapi.InvalidParameter
		}
		condition := form.Condition{
			MetricName:         monitorItem.Name,
			MetricCode:         param.MetricCode,
			Period:             cond.Period,
			Times:              cond.Times,
			Statistics:         cond.Statistics,
			ComparisonOperator: cond.ComparisonOperator,
			Threshold:          util.FormatFloat(threshold, 2),
			Unit:               monitorItem.Unit,
			Labels:             monitorItem.Labels,
			Level:              cond.Level,
			SilencesTime:       "3h",
		}
		if condition.Level == 0 && param.AlarmLevel > 0 {
			condition.Level = param.AlarmLevel
		}
		if len(condition.SilencesTime) == 0 && len(param.SilencesTime) > 0 {
			condition.SilencesTime = param.SilencesTime
		}
		addForm.Conditions = append(addForm.Conditions, condition)
	}
	return nil
}

func buildMultipleAlarmRuleReqParam(param *AlarmRuleCreateReqDTO, addForm *form.AlarmRuleAddReqDTO) *openapi.ErrorCode {
	addForm.Combination = param.Combination
	for _, cond := range param.TriggerConditions {
		monitorItem, err1 := checkMetricName(cond.MetricCode)
		if err1 != nil {
			return openapi.MetricCodeInvalid
		}
		matched, err := regexp.MatchString("^[0-9\\.]+$", cond.Threshold)
		if !matched {
			return openapi.ThresholdInvalid
		}
		threshold, err := strconv.ParseFloat(cond.Threshold, 64)
		if err != nil {
			logger.Logger().Infof("Threshold is parsefloat error:%v", err)
			return openapi.InvalidParameter
		}
		addForm.Period = param.Period
		addForm.Times = param.Times
		condition := form.Condition{
			MetricName:         monitorItem.Name,
			MetricCode:         cond.MetricCode,
			Period:             param.Period,
			Times:              param.Times,
			Statistics:         cond.Statistics,
			ComparisonOperator: cond.ComparisonOperator,
			Threshold:          util.FormatFloat(threshold, 2),
			Unit:               monitorItem.Unit,
			Labels:             monitorItem.Labels,
			Level:              param.AlarmLevel,
			SilencesTime:       "3h",
		}
		if condition.Level == 0 && param.AlarmLevel > 0 {
			condition.Level = param.AlarmLevel
		}
		if len(condition.SilencesTime) == 0 && len(param.SilencesTime) > 0 {
			condition.SilencesTime = param.SilencesTime
		}
		addForm.Conditions = append(addForm.Conditions, condition)
	}
	return nil
}

func CheckAndFillParam(param *form.AlarmRuleAddReqDTO) *openapi.ErrorCode {
	if param.Type == constant.AlarmRuleTypeSingleMetric {
		if len(param.MetricCode) == 0 {
			return openapi.MetricCodeMissing
		}
		if len(param.SilencesTime) == 0 {
			return openapi.RuleSilencesTimeMissing
		}
		for i, cond := range param.Conditions {
			param.Conditions[i].MetricCode = param.MetricCode
			param.Conditions[i].SilencesTime = param.SilencesTime
			if len(cond.MetricCode) == 0 {
				return openapi.MetricCodeMissing
			}
			if cond.Level == 0 {
				return openapi.RuleLevelMissing
			}
			if len(cond.SilencesTime) == 0 {
				return openapi.RuleSilencesTimeMissing
			}
			if cond.Period == 0 {
				return openapi.RulePeriodMissing
			}
			if cond.Times == 0 {
				return openapi.RuleTimesMissing
			}
		}
	} else {
		if param.Level == 0 {
			return openapi.RuleLevelMissing
		}
		if len(param.SilencesTime) == 0 {
			return openapi.RuleSilencesTimeMissing
		}
		if param.Period == 0 {
			return openapi.RulePeriodMissing
		}
		if param.Times == 0 {
			return openapi.RuleTimesMissing
		}
		if param.Combination != 1 && param.Combination != 2 {
			return openapi.RuleCombinationMissing
		}
	}
	return nil
}

type AlarmPageReqParam struct {
	RuleName   string
	Status     string `binding:"oneof=disabled enabled '' "`
	PageNumber int    `binding:"min=1"`
	PageSize   int    `binding:"min=1,max=100"`
}

type AlarmRulePageDTO struct {
	openapi.ResCommonPage
	Rules []RuleInfo
}
type RuleInfo struct {
	Name                string
	MonitorType         string
	ProductType         string
	ResourceNum         int
	Status              string
	RuleId              string
	ProductAbbreviation string
}

type RuleDetail struct {
	RequestId string
	RuleInfo  struct {
		Status         string
		RuleId         string
		Type           uint8
		RuleConditions []struct {
			MetricCode         string
			MetricName         string
			Period             int
			Times              int
			Statistics         string
			ComparisonOperator string
			Threshold          string
			Unit               string
			Labels             string
			Level              uint8
			SilencesTime       string
			Express            string
		}
		NoticeGroups []struct {
			Id   string
			Name string
		}
		MonitorType  string
		ProductType  string
		Scope        string
		RuleName     string
		SilencesTime string
		AlarmLevel   int
		Combination  uint8
		Describe     string
		ResourceIds  []struct {
			ResourceId string
		}
		NoticeChannels []struct {
			HandlerType int
		}
		Period int
		Times  int
	}
}

type AlarmRuleCreateReqDTO struct {
	Type                uint8  `binding:"oneof=1 2"`
	MonitorType         string `binding:"oneof=云产品监控"`
	ProductAbbreviation string `binding:"required"`
	Scope               string `binding:"oneof=ALL Resource"`
	Resources           []struct {
		ResourceId string `binding:"required"`
	} `binding:"required"`
	NoticeChannels []struct {
		HandlerType int `binding:"oneof=1 2"`
	}
	RuleName          string `binding:"required"`
	TriggerConditions []struct {
		MetricCode         string
		Period             int    `binding:"required"`
		Times              int    `binding:"required"`
		Statistics         string `binding:"oneof=Maximum Minimum Average"`
		ComparisonOperator string `binding:"oneof=greater greaterOrEqual less  lessOrEqual  equal notEqual"`
		Threshold          string
		//TODO 枚举
		Level        uint8
		SilencesTime string
	} `binding:"required"`
	SilencesTime string
	AlarmLevel   uint8 ` binding:"oneof=0 1 2 3 4"`
	GroupList    []string
	MetricCode   string
	Period       int
	Times        int
	Combination  uint8
}

type AlarmRuleUpdateReqDTO struct {
	Type      uint8  `binding:"oneof=1 2"`
	Scope     string `binding:"oneof=ALL Resource"`
	Resources []struct {
		ResourceId string `binding:"required"`
	} `binding:"required"`
	NoticeChannels []struct {
		HandlerType int `binding:"oneof=1 2"`
	}
	RuleName          string `binding:"required"`
	TriggerConditions []struct {
		MetricCode         string
		Period             int    `binding:"required"`
		Times              int    `binding:"required"`
		Statistics         string `binding:"oneof=Maximum Minimum Average"`
		ComparisonOperator string `binding:"oneof=greater greaterOrEqual less  lessOrEqual  equal notEqual"`
		Threshold          string
		//TODO 枚举
		Level        uint8
		SilencesTime string
	} `binding:"required"`
	SilencesTime string
	AlarmLevel   uint8 ` binding:"oneof=0 1 2 3 4 "`
	GroupList    []string
	MetricCode   string
	Period       int
	Times        int
	Combination  uint8
}

type StatusBody struct {
	Enable string `binding:"required,oneof=disabled enabled"`
}
