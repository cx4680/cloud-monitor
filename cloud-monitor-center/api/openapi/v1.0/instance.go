package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonError "code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceCtl struct{}

func NewInstanceCtl() *InstanceCtl {
	return &InstanceCtl{}
}

func (ctl *InstanceCtl) Page(c *gin.Context) {
	reqParam := openapi.NewPageQuery()
	if err := c.ShouldBindQuery(&reqParam); err != nil {
		logger.Logger().Infof("param valid error %+v", err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	resourceId := c.Param("ResourceId")
	pageParam := form.InstanceRulePageReqParam{
		InstanceId: resourceId,
		PageSize:   reqParam.PageSize,
		Current:    reqParam.PageNumber,
	}
	page := dao.Instance.SelectInstanceRulePage(&pageParam)
	pageVO := page.(*vo.PageVO)
	listVo := pageVO.Records.([]*form.InstanceRuleDTO)
	var ruleList []InstanceRuleInfo
	for _, item := range listVo {
		ruleInfo := InstanceRuleInfo{
			MonitorType: item.MonitorType,
			RuleName:    item.Name,
			RuleId:      item.Id,
			MetricCode:  item.RuleCondition.MetricName,
			MetricName:  item.RuleCondition.MonitorItemName,
		}
		ruleList = append(ruleList, ruleInfo)
	}
	rulePage := InstanceRulePage{
		ResCommonPage: *openapi.NewResCommonPage(c, pageVO),
		Rules:         ruleList,
	}
	c.JSON(http.StatusOK, rulePage)
}

func (ctl *InstanceCtl) Unbind(c *gin.Context) {
	var reqParam UnBindBodyParam
	if err := c.ShouldBindJSON(&reqParam); err != nil {
		logger.Logger().Infof("param valid error %+v", err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	tenantId, err2 := util.GetTenantId(c)
	if err2 != nil {
		logger.Logger().Info("tenantId is nil")
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	resourceId := c.Param("ResourceId")
	unBindParam := form.UnBindRuleParam{
		InstanceId: resourceId,
		RuleId:     reqParam.RuleId,
		TenantId:   tenantId,
	}
	err := util.Tx(&unBindParam, service.UnbindInstance)
	if err != nil {
		logger.Logger().Info(err)
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.RuleIdInvalid, c))
			return
		}
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

func (ctl *InstanceCtl) Bind(c *gin.Context) {
	var reqParam BindBodyParam
	if err := c.ShouldBindJSON(&reqParam); err != nil {
		logger.Logger().Infof("param valid error %+v", err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		logger.Logger().Info("tenantId is nil")
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
		return
	}
	resourceId := c.Param("ResourceId")
	instanceParam := form.InstanceBindRuleDTO{
		TenantId: tenantId,
		InstanceInfo: form.InstanceInfo{
			InstanceId: resourceId,
		},
		RuleIdList: reqParam.RuleIds,
	}
	err = util.Tx(&instanceParam, service.BindInstance)
	if err != nil {
		logger.Logger().Info(err)
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.InvalidParameter, c))
			return
		}
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	c.JSON(http.StatusOK, openapi.NewResSuccess(c))
}

type InstanceRuleInfo struct {
	MonitorType string
	RuleName    string
	RuleId      string
	MetricName  string
	MetricCode  string
}

type InstanceRulePage struct {
	openapi.ResCommonPage
	Rules []InstanceRuleInfo
}

type BindBodyParam struct {
	RuleIds []string
}

type UnBindBodyParam struct {
	RuleId string `binding:"required"`
}
