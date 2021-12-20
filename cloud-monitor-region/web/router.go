package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/actuator"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/controllers"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/docs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/inner"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func loadRouters() {
	monitorProductRouters()
	swagger()
	instance()
	MonitorReportForm()
	innerCtl()
	actuatorMapping()
}

func MonitorReportForm() {
	monitorReportFormCtl := controllers.NewMonitorReportFormController(service.NewMonitorReportFormService())
	group := router.Group("/hawkeye/MonitorReportForm/")
	{
		group.GET("/getData", monitorReportFormCtl.GetData)
		group.GET("/getAxisData", monitorReportFormCtl.GetAxisData)
		group.GET("/getTop", monitorReportFormCtl.GetTop)
	}
}

func monitorProductRouters() {
	monitorProductCtl := controllers.NewMonitorProductCtl(commonDao.MonitorProduct)
	group := router.Group("/hawkeye/monitorProduct/")
	{
		group.GET("/getById", monitorProductCtl.GetById)
		group.PUT("/updateById", monitorProductCtl.UpdateById)
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
		group.GET("/page", instanceCtl.GetPage)
		group.GET("/getInstanceNum", instanceCtl.GetInstanceNumByRegion)
	}
}

func innerCtl() {
	addService := service.NewAlertRecordAddService(service.NewAlertRecordService(), commonService.NewAlarmHandlerService(), commonService.NewMessageService(messageCenter.NewService()), commonService.NewTenantService())
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
