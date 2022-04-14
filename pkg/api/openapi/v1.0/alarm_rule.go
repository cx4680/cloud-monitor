package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	commonError "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	form "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	util2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
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
	tenantId, err := util2.GetTenantId(c)
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
				MetricCode:          ruleVo.MetricName,
				MetricName:          ruleVo.RuleCondition.MonitorItemName,
				Express:             ruleVo.Express,
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
	tenantId, err := util2.GetTenantId(c)
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
	var scope string
	if detailVo.Scope == "INSTANCE" {
		scope = "Resource"
	}
	res.RuleInfo.Status = detailVo.Status
	res.RuleInfo.Scope = scope
	res.RuleInfo.MonitorType = detailVo.MonitorType
	res.RuleInfo.RuleName = detailVo.RuleName
	res.RuleInfo.SilencesTime = detailVo.SilencesTime
	res.RuleInfo.AlarmLevel = detailVo.AlarmLevel
	res.RuleInfo.RuleCondition.Unit = detailVo.RuleCondition.Unit
	res.RuleInfo.RuleCondition.MetricCode = detailVo.RuleCondition.MetricName
	res.RuleInfo.RuleCondition.MetricName = detailVo.RuleCondition.MonitorItemName
	res.RuleInfo.RuleCondition.Period = strconv.Itoa(detailVo.RuleCondition.Period)
	res.RuleInfo.RuleCondition.Times = detailVo.RuleCondition.Times
	res.RuleInfo.RuleCondition.Statistics = detailVo.RuleCondition.Statistics
	res.RuleInfo.RuleCondition.ComparisonOperator = detailVo.RuleCondition.ComparisonOperator
	res.RuleInfo.RuleCondition.Threshold = fmt.Sprintf("%.2f", detailVo.RuleCondition.Threshold)
	res.RuleInfo.RuleCondition.Labels = detailVo.RuleCondition.Labels
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
	err := util2.Tx(addForm, service.CreateRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	res := struct {
		RequestId string
		RuleId    string
	}{
		RequestId: openapi.GetRequestId(c),
		RuleId:    addForm.Id,
	}
	c.JSON(http.StatusOK, res)
}

func buildAlarmRuleReqParam(c *gin.Context, createParam *AlarmRuleCreateReqDTO, updateParam *AlarmRuleUpdateReqDTO) (*form.AlarmRuleAddReqDTO, *openapi.ErrorCode) {
	var param AlarmRuleCreateReqDTO
	var productInfo = &model.MonitorProduct{}
	if createParam != nil {
		if err := c.ShouldBindJSON(&createParam); err != nil {
			logger.Logger().Info(err)
			return nil, openapi.GetErrorCode(err)
		}
		param = *createParam
		productInfo = dao.MonitorProduct.GetByAbbreviation(global.DB, param.ProductAbbreviation)
		if productInfo == nil || len(productInfo.BizId) == 0 {
			return nil, openapi.ProductAbbreviationInvalid
		}
	} else if updateParam != nil {
		if err := c.ShouldBindJSON(&updateParam); err != nil {
			logger.Logger().Info(err)
			return nil, openapi.GetErrorCode(err)
		}
		param = AlarmRuleCreateReqDTO{
			Scope:            updateParam.Scope,
			Resources:        updateParam.Resources,
			RuleName:         updateParam.RuleName,
			NoticeChannels:   updateParam.NoticeChannels,
			TriggerCondition: updateParam.TriggerCondition,
			SilencesTime:     updateParam.SilencesTime,
			AlarmLevel:       updateParam.AlarmLevel,
			GroupList:        updateParam.GroupList,
		}
	}
	nameMatched, err := regexp.MatchString("^[a-z][a-z0-9_]{0,14}[a-z0-9]$", param.RuleName)
	if !nameMatched {
		return nil, openapi.RuleNameInvalid
	}

	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		logger.Logger().Info("tenantId is nil")
		return nil, openapi.MissingParameter
	}
	userId, err := util2.GetUserId(c)
	if err != nil {
		logger.Logger().Info("userId is nil")
		return nil, openapi.MissingParameter
	}
	monitorItem, err1 := checkMetricName(param.TriggerCondition.MetricCode)
	if err1 != nil {
		return nil, openapi.MetricCodeInvalid
	}
	productBizId, _ := strconv.Atoi(productInfo.BizId)
	var scope string
	if param.Scope == "Resource" {
		scope = "INSTANCE"
	}
	addForm := form.AlarmRuleAddReqDTO{
		RuleName:     param.RuleName,
		MonitorType:  "云产品监控",
		ProductType:  productInfo.Name,
		ProductId:    productBizId,
		Scope:        scope,
		TenantId:     tenantId,
		UserId:       userId,
		SilencesTime: "3h",
		AlarmLevel:   param.AlarmLevel,
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
	matched, err := regexp.MatchString("^[0-9\\.]+$", param.TriggerCondition.Threshold)
	if !matched {
		return nil, openapi.ThresholdInvalid
	}
	threshold, err := strconv.ParseFloat(param.TriggerCondition.Threshold, 64)
	if err != nil {
		logger.Logger().Infof("Threshold is parsefloat error:%v", err)
		return nil, openapi.InvalidParameter
	}
	addForm.RuleCondition = &form.RuleCondition{
		MetricName:         param.TriggerCondition.MetricCode,
		Period:             param.TriggerCondition.Period,
		Times:              param.TriggerCondition.Times,
		Statistics:         param.TriggerCondition.Statistics,
		ComparisonOperator: param.TriggerCondition.ComparisonOperator,
		Threshold:          util2.FormatFloat(threshold, 2),
		Unit:               monitorItem.Unit,
		Labels:             monitorItem.Labels,
		MonitorItemName:    monitorItem.Name,
	}
	return &addForm, nil
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var reqParam AlarmRuleUpdateReqDTO
	updateForm, errCode := buildAlarmRuleReqParam(c, nil, &reqParam)
	if errCode != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(errCode, c))
		return
	}
	ruleId := c.Param("RuleId")
	updateForm.Id = ruleId
	err := util2.Tx(updateForm, service.UpdateRule)
	if err != nil {
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
			return
		}
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func (ctl *AlarmRuleCtl) DeleteRule(c *gin.Context) {
	ruleId := c.Param("RuleId")
	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	reqParam := form.RuleReqDTO{
		Id:       ruleId,
		TenantId: tenantId,
	}
	err = util2.Tx(&reqParam, service.DeleteRule)
	if err != nil {
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
			return
		}
		c.JSON(http.StatusInternalServerError, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func (ctl *AlarmRuleCtl) ChangeRuleStatus(c *gin.Context) {
	ruleId := c.Param("RuleId")
	var status StatusBody
	if err := c.ShouldBindJSON(&status); err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	tenantId, err := util2.GetTenantId(c)
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
	err = util2.Tx(&reqParam, service.ChangeRuleStatus)
	if err != nil {
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
			return
		}
		logger.Logger().Info(err)
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func checkMetricName(metricCode string) (*model.MonitorItem, error) {
	item := dao.AlarmRule.GetMonitorItem(metricCode)
	if item == nil || len(item.MetricName) == 0 {
		return nil, errors.New("指标不存在")
	}
	return item, nil
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
	MetricCode          string
	MetricName          string
	Express             string
	ResourceNum         int
	Status              string
	RuleId              string
	ProductAbbreviation string
}

type RuleDetail struct {
	RequestId string
	RuleInfo  struct {
		Status        string
		RuleCondition struct {
			MetricCode         string
			MetricName         string
			Period             string
			Times              int
			Statistics         string
			ComparisonOperator string
			Threshold          string
			Unit               string
			Labels             string
		}
		NoticeGroups []struct {
			Id   string
			Name string
		}
		MonitorType  string
		Scope        string
		RuleName     string
		SilencesTime string
		AlarmLevel   int
		ResourceIds  []struct {
			ResourceId string
		}
		NoticeChannels []struct {
			HandlerType int
		}
	}
}

type AlarmRuleCreateReqDTO struct {
	MonitorType         string `binding:"oneof=云产品监控"`
	ProductAbbreviation string `binding:"required"`
	Scope               string `binding:"oneof=ALL Resource"`
	Resources           []struct {
		ResourceId string `binding:"required"`
	} `binding:"required"`
	NoticeChannels []struct {
		HandlerType int `binding:"oneof=1 2"`
	}
	RuleName         string `binding:"required"`
	TriggerCondition struct {
		MetricCode         string `binding:"required"`
		Period             int    `binding:"required"`
		Times              int    `binding:"required"`
		Statistics         string `binding:"oneof=Maximum Minimum Average"`
		ComparisonOperator string `binding:"oneof=greater greaterOrEqual less  lessOrEqual  equal notEqual"`
		Threshold          string
	} `binding:"required"`
	SilencesTime string
	AlarmLevel   uint8 ` binding:"oneof=1  2 3 4 "`
	GroupList    []string
}

type AlarmRuleUpdateReqDTO struct {
	Scope     string `binding:"oneof=ALL Resource"`
	Resources []struct {
		ResourceId string `binding:"required"`
	} `binding:"required"`
	NoticeChannels []struct {
		HandlerType int `binding:"oneof=1 2"`
	}
	RuleName         string `binding:"required"`
	TriggerCondition struct {
		MetricCode         string `binding:"required"`
		Period             int    `binding:"required"`
		Times              int    `binding:"required"`
		Statistics         string `binding:"oneof=Maximum Minimum Average"`
		ComparisonOperator string `binding:"oneof=greater greaterOrEqual less  lessOrEqual  equal notEqual"`
		Threshold          string
	} `binding:"required"`
	SilencesTime string
	AlarmLevel   uint8 ` binding:"oneof=1  2 3 4 "`
	GroupList    []string
}

type StatusBody struct {
	Enable string `binding:"required,oneof=disabled enabled"`
}
