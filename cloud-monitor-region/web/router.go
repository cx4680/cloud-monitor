package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
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
		group.GET("/getData", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetData)
		group.GET("/getAxisData", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportRangeData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetAxisData)
		group.GET("/getTop", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportTop", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetTop)
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
		group.GET("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstancePageList", ResourceType: "*", ResourceId: "*"}), instanceCtl.GetPage)
		group.GET("/getInstanceNum", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceNum", ResourceType: "*", ResourceId: "*"}), instanceCtl.GetInstanceNumByRegion)
	}
}

func innerCtl() {
	addService := service.NewAlertRecordAddService(service.NewAlertRecordService(), commonService.NewAlarmHandlerService(), commonService.NewMessageService(message_center.NewService()), commonService.NewTenantService())
	ctl := inner.NewAlertRecordCtl(addService)
	group := router.Group("/inner/")
	{
		group.POST("/alertRecord/insert", ctl.AddAlertRecord)
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
