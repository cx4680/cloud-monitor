package v1_0

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	commonUtil "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/external"
	_ "code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResourceCtl struct {
}

func NewResourceController() *ResourceCtl {
	return &ResourceCtl{}
}
func (ctl *ResourceCtl) GetResourceList(c *gin.Context) {
	tenantId, err := commonUtil.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.MissingParameter, c))
	}
	param := openapi.NewPageQuery()
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	productCode := c.Param("ProductCode")
	c.Set(global.ResourceName, productCode)
	f := commonService.InstancePageForm{Product: productCode, TenantId: tenantId, PageSize: param.PageSize, Current: param.PageNumber}
	instanceService := external.ProductInstanceServiceMap[f.Product]
	if instanceService == nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.ProductCodeInvalid, c))
		return
	}
	page, err := instanceService.GetPage(f, instanceService.(commonService.InstanceStage))
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
		return
	}
	var resources []ResourceInfo
	if records, ok := page.Records.([]commonService.InstanceCommonVO); ok {
		for _, record := range records {
			resources = append(resources, ResourceInfo{
				ResourceId:   record.InstanceId,
				ResourceName: record.InstanceName,
			})
		}
	}

	resourcePage := ResourcePage{
		ResCommonPage: &openapi.ResCommonPage{
			ResCommon: openapi.ResCommon{
				RequestId: openapi.GetRequestId(c),
			},
			TotalCount: page.Total,
			PageSize:   page.Size,
			PageNumber: page.Current,
		},
		Resources: resources,
	}
	c.JSON(http.StatusOK, resourcePage)
}

type ResourcePage struct {
	*openapi.ResCommonPage
	Resources []ResourceInfo
}

type ResourceInfo struct {
	ResourceId   string
	ResourceName string
}
