package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/controllers"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/docs"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func loadRouters() {
	monitorProductRouters()
	swagger()
	instance()
	eip()
	slb()
}

func monitorProductRouters() {
	monitorProductCtl := controllers.NewMonitorProductCtl(dao.NewMonitorProductDao(database.GetDb()))
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
	instanceCtl := controllers.NewInstanceCtl(dao.NewInstanceDao(database.GetDb()))
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
