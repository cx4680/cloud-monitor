package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/logs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/api/openapi/v1.0"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
)

func loadOpenApiV1Routers() {
	group := router.Group("/v1.0/")
	MonitorReportOpenApiV1Routers(group)
	ResourceOpenApiV1Routers(group)
}

func MonitorReportOpenApiV1Routers(group *gin.RouterGroup) {
	monitorReportFormCtl := v1_0.NewMonitorReportFormController()
	group.GET("resources/:ResourceId/metrics/:MetricCode/datas", logs.GinTrailzap(false, Read, logs.INFO, logs.MonitorReportForm), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportRangeData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetMonitorDatas)
	group.GET("resources/:ResourceId/metrics/:MetricCode/data", logs.GinTrailzap(false, Read, logs.INFO, logs.MonitorReportForm), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetMonitorData)
	group.GET("metrics/:MetricCode/:N/resources", logs.GinTrailzap(false, Read, logs.INFO, logs.MonitorReportForm), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportTop", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetMonitorDataTop)
}

func ResourceOpenApiV1Routers(group *gin.RouterGroup) {
	resourceCtl := v1_0.NewResourceController()
	group.GET(":ProductAbbreviation/resources", logs.GinTrailzap(false, Read, logs.INFO, logs.Resource), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetResourceList", ResourceType: "*", ResourceId: "*"}), resourceCtl.GetResourceList)
}
