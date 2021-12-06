package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/controllers"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/docs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/inner"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func loadRouters() {
	monitorProductRouters()
	swagger()
	instance()
	eip()
	slb()
	nat()
	MonitorReportForm()
	innerCtl()
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
		group.GET("/page", instanceCtl.Page)
		group.GET("/getInstanceNum", instanceCtl.GetInstanceNumByRegion)
	}
}

func eip() {
	ctl := controllers.NewEipCtl()
	group := router.Group("/hawkeye/eip/")
	{
		group.GET("/page", ctl.Page)
	}
}

func slb() {
	ctl := controllers.NewSlbCtl()
	group := router.Group("/hawkeye/slb/")
	{
		group.GET("/page", ctl.Page)
	}
}

func nat() {
	ctl := controllers.NewNatCtl()
	group := router.Group("/hawkeye/nat/")
	{
		group.GET("/page", ctl.Page)
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
