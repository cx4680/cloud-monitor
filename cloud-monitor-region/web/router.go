package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/logs"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/api/actuator"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/api/controller"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/api/inner"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/docs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

const Read = "Read"

func loadRouters() {
	swagger()
	instance()
	MonitorReportForm()
	innerCtl()
	actuatorMapping()
	remote()
}

func MonitorReportForm() {
	monitorReportFormCtl := controller.NewMonitorReportFormController(service.NewMonitorReportFormService())
	group := router.Group("/hawkeye/MonitorReportForm/")
	{
		group.GET("/getData", logs.GinTrailzap(false, Read, logs.INFO, logs.MonitorReportForm), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetData)
		group.GET("/getAxisData", logs.GinTrailzap(false, Read, logs.INFO, logs.MonitorReportForm), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportRangeData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetAxisData)
		group.GET("/getTop", logs.GinTrailzap(false, Read, logs.INFO, logs.MonitorReportForm), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportTop", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetTop)
	}
}

func swagger() {
	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
}

func instance() {
	instanceCtl := controller.NewInstanceCtl(dao.Instance)
	group := router.Group("/hawkeye/instance/")
	{
		group.GET("/page", logs.GinTrailzap(false, Read, logs.INFO, logs.Resource), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstancePageList", ResourceType: "*", ResourceId: "*"}), instanceCtl.GetPage)
		group.GET("/getInstanceNum", logs.GinTrailzap(false, Read, logs.INFO, logs.Resource), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceNum", ResourceType: "*", ResourceId: "*"}), instanceCtl.GetInstanceNumByRegion)
	}
}

func innerCtl() {
	addService := service.NewAlarmRecordAddService(service.NewAlarmRecordService(service.NewAlarmInfoService()), commonService.NewAlarmHandlerService(), commonService.NewTenantService())
	ctl := inner.NewAlertRecordCtl(addService)
	group := router.Group("/inner/")
	{
		group.POST("/alarmRecord/insert", ctl.AddAlarmRecord)
	}
}

func actuatorMapping() {
	group := router.Group("/actuator")
	{
		group.GET("/env", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator.Env())
		})
		group.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator.Info())
		})
		group.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator.Health())
		})
		group.GET("/metrics", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator.Metrics())
		})
	}
}

func remote() {
	router.GET("/inner/remote/:productType", func(context *gin.Context) {
		productType := context.Param("productType")
		task.Run(productType)
	})
}
