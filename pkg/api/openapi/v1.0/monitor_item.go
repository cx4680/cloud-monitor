package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorItemCtl struct {
	service service.MonitorItemService
}

func NewMonitorItemCtl(service service.MonitorItemService) *MonitorItemCtl {
	return &MonitorItemCtl{service}
}

func (mic *MonitorItemCtl) GetMonitorItemsByProductAbbr(c *gin.Context) {
	param := openapi.NewPageQuery()
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	productAbbreviation := c.Param("ProductAbbreviation")
	c.Set(global.ResourceName, productAbbreviation)
	abbreviation := dao.MonitorProduct.GetByAbbreviation(global.DB, productAbbreviation)
	if len(abbreviation.BizId) == 0 {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.ProductAbbreviationInvalid, c))
		return
	}
	c.Set(global.ResourceName, productAbbreviation)
	pageVo := mic.service.GetMonitorItemPage(param.PageSize, param.PageNumber, productAbbreviation)
	var metricMetaList []MetricMeta
	metricListVo := *pageVo.Records.(*[]model.MonitorItem)
	for _, metricVo := range metricListVo {
		metricMeta := MetricMeta{
			Name:                metricVo.Name,
			Unit:                metricVo.Unit,
			Description:         metricVo.Description,
			MetricCode:          metricVo.MetricName,
			ProductAbbreviation: productAbbreviation,
			Dimensions:          metricVo.Labels,
		}
		metricMetaList = append(metricMetaList, metricMeta)
	}
	metricPage := MetricPage{
		ResCommonPage: *openapi.NewResCommonPage(c, pageVo),
		Metrics:       metricMetaList,
	}
	c.JSON(http.StatusOK, metricPage)
}

type MetricMeta struct {
	Name                string
	Unit                string
	ProductAbbreviation string
	MetricCode          string
	Description         string
	Dimensions          string
}

type MetricPage struct {
	openapi.ResCommonPage
	Metrics []MetricMeta
}
