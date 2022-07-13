package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorProductCtl struct {
	service service.MonitorProductService
}

func NewMonitorProductCtl(service service.MonitorProductService) *MonitorProductCtl {
	return &MonitorProductCtl{service}
}

func (mpc *MonitorProductCtl) GetMonitorProduct(c *gin.Context) {
	pageQuery := openapi.NewPageQuery()
	if err := c.ShouldBindQuery(pageQuery); err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	productPageVo := mpc.service.GetMonitorProductPage(pageQuery.PageSize, pageQuery.PageNumber)
	var productMetaList []ProductMeta
	productListVo := productPageVo.Records.([]model.MonitorProduct)
	for _, productVo := range productListVo {
		productMeta := ProductMeta{
			ProductName: productVo.Name,
			ProductCode: productVo.Abbreviation,
			Description: productVo.Description,
			MonitorType: productVo.MonitorType,
		}
		productMetaList = append(productMetaList, productMeta)
	}
	page := ProductPage{
		ResCommonPage: *openapi.NewResCommonPage(c, productPageVo),
		ProductList:   productMetaList,
	}
	c.JSON(http.StatusOK, page)
}

type ProductMeta struct {
	ProductName string // 监控产品名称
	Description string // 描述
	ProductCode string // 监控产品编码
	MonitorType string // 监控类型
}

type ProductPage struct {
	openapi.ResCommonPage
	ProductList []ProductMeta
}
