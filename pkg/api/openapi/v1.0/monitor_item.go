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

func (mic *MonitorItemCtl) GetMonitorItemsByProductCode(c *gin.Context) {
	param := openapi.NewPageQuery()
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	ProductCode := c.Param("ProductCode")
	c.Set(global.ResourceName, ProductCode)
	monitorProduct := dao.MonitorProduct.GetByProductCode(global.DB, ProductCode)
	if len(monitorProduct.BizId) == 0 {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.ProductCodeInvalid, c))
		return
	}
	c.Set(global.ResourceName, ProductCode)
	pageVo := mic.service.GetMonitorItemPage(param.PageSize, param.PageNumber, ProductCode)
	var metricMetaList []MetricMeta
	metricListVo := *pageVo.Records.(*[]model.MonitorItem)
	for _, metricVo := range metricListVo {
		metricMeta := MetricMeta{
			Name:        metricVo.Name,
			Unit:        metricVo.Unit,
			Description: metricVo.Description,
			MetricCode:  metricVo.MetricName,
			Period:      []int{5, 15, 30, 60, 3600},
			Periods:     []Periods{{5, "AVG,MIN,MAX"}, {15, "AVG,MIN,MAX"}, {30, "AVG,MIN,MAX"}, {60, "AVG,MIN,MAX"}, {3600, "AVG,MIN,MAX"}},
			ProductCode: ProductCode,
			Dimensions:  metricVo.Labels,
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
	Name        string
	Unit        string
	ProductCode string
	MetricCode  string
	Period      []int
	Periods     []Periods
	Description string
	Dimensions  string
}

type Periods struct {
	Period   int
	StatType string
}

type MetricPage struct {
	openapi.ResCommonPage
	Metrics []MetricMeta
}
