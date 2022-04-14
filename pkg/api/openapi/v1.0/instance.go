package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	commonError "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	openapi2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	util2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceCtl struct{}

func NewInstanceCtl() *InstanceCtl {
	return &InstanceCtl{}
}

func (ctl *InstanceCtl) Page(c *gin.Context) {
	reqParam := openapi2.NewPageQuery()
	if err := c.ShouldBindQuery(&reqParam); err != nil {
		logger.Logger().Infof("param valid error %+v", err)
		c.JSON(http.StatusBadRequest, openapi2.NewRespError(openapi2.GetErrorCode(err), c))
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
		ResCommonPage: *openapi2.NewResCommonPage(c, pageVO),
		Rules:         ruleList,
	}
	c.JSON(http.StatusOK, rulePage)
}

func (ctl *InstanceCtl) Unbind(c *gin.Context) {
	var reqParam UnBindBodyParam
	if err := c.ShouldBindJSON(&reqParam); err != nil {
		logger.Logger().Infof("param valid error %+v", err)
		c.JSON(http.StatusBadRequest, openapi2.NewRespError(openapi2.GetErrorCode(err), c))
		return
	}
	tenantId, err2 := util2.GetTenantId(c)
	if err2 != nil {
		logger.Logger().Info("tenantId is nil")
		c.JSON(http.StatusBadRequest, openapi2.NewRespError(openapi2.MissingParameter, c))
		return
	}
	resourceId := c.Param("ResourceId")
	unBindParam := form.UnBindRuleParam{
		InstanceId: resourceId,
		RuleId:     reqParam.RuleId,
		TenantId:   tenantId,
	}
	err := util2.Tx(&unBindParam, service.UnbindInstance)
	if err != nil {
		logger.Logger().Info(err)
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi2.NewRespError(openapi2.RuleIdInvalid, c))
			return
		}
		c.JSON(http.StatusInternalServerError, openapi2.NewRespError(openapi2.SystemError, c))
		return
	}
	c.JSON(http.StatusOK, openapi2.NewResSuccess(c))
}

func (ctl *InstanceCtl) Bind(c *gin.Context) {
	var reqParam BindBodyParam
	if err := c.ShouldBindJSON(&reqParam); err != nil {
		logger.Logger().Infof("param valid error %+v", err)
		c.JSON(http.StatusBadRequest, openapi2.NewRespError(openapi2.GetErrorCode(err), c))
		return
	}
	tenantId, err := util2.GetTenantId(c)
	if err != nil {
		logger.Logger().Info("tenantId is nil")
		c.JSON(http.StatusBadRequest, openapi2.NewRespError(openapi2.MissingParameter, c))
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
	err = util2.Tx(&instanceParam, service.BindInstance)
	if err != nil {
		logger.Logger().Info(err)
		if _, ok := err.(*commonError.BusinessError); ok {
			c.JSON(http.StatusBadRequest, openapi2.NewRespError(openapi2.InvalidParameter, c))
			return
		}
		c.JSON(http.StatusInternalServerError, openapi2.NewRespError(openapi2.SystemError, c))
		return
	}
	c.JSON(http.StatusOK, openapi2.NewResSuccess(c))
}

type InstanceRuleInfo struct {
	MonitorType string
	RuleName    string
	RuleId      string
	MetricName  string
	MetricCode  string
}

type InstanceRulePage struct {
	openapi2.ResCommonPage
	Rules []InstanceRuleInfo
}

type BindBodyParam struct {
	RuleIds []string
}

type UnBindBodyParam struct {
	RuleId string `binding:"required"`
}
