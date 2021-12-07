package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/actuator"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/controllers"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/docs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/inner"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func loadRouters() {
	swagger()
	instance()
	eip()
	slb()
	nat()
	cbr()
	MonitorReportForm()
	innerCtl()
	actuatorMapping()
}

func MonitorReportForm() {
	monitorReportFormCtl := controllers.NewMonitorReportFormController(service.NewMonitorReportFormService())
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
	instanceCtl := controllers.NewInstanceCtl(dao.Instance)
	group := router.Group("/hawkeye/instance/")
	{
		group.GET("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstancePageList", ResourceType: "*", ResourceId: "*"}), instanceCtl.Page)
		group.GET("/getInstanceNum", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceNum", ResourceType: "*", ResourceId: "*"}), instanceCtl.GetInstanceNumByRegion)
	}
}

func eip() {
	ctl := controllers.NewEipCtl()
	group := router.Group("/hawkeye/eip/")
	{
		group.GET("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetEipPageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
	}
}

func slb() {
	ctl := controllers.NewSlbCtl()
	group := router.Group("/hawkeye/slb/")
	{
		group.GET("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetSlbPageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
	}
}

func cbr() {
	ctl := controllers.NewCbrCtl()
	group := router.Group("/hawkeye/cbr/")
	{
		group.GET("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetCbrPageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
	}
}

func nat() {
	ctl := controllers.NewNatCtl()
	group := router.Group("/hawkeye/nat/")
	{
		group.GET("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNatPageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
	}
}

func innerCtl() {
	addService := service.NewAlertRecordAddService(service.NewAlertRecordService(), commonService.NewMessageService(messageCenter.NewService()), commonService.NewTenantService())
	ctl := inner.NewAlertRecordCtl(addService)
	group := router.Group("/hawkeye/inner/")
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
