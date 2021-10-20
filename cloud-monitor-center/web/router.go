package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/controllers"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/database"
)

func loadRouters() {
	monitorProductRouters()
	alertContactRouters()
	alertContactGroupRouters()
}

func monitorProductRouters() {
	monitorProductCtl := controllers.NewMonitorProductCtl(dao.NewMonitorProductDao(database.GetDb()))
	group := router.Group("/hawkeye/monitorProduct/")
	{
		group.GET("/getById", monitorProductCtl.GetById)
		group.PUT("/updateById", monitorProductCtl.UpdateById)
	}
}

func alertContactRouters() {
	alertContactCtl := controllers.NewAlertContactCtl(dao.NewAlertContact(database.GetDb()))
	group := router.Group("/hawkeye/alertContact/")
	{
		group.GET("/getAlertContact", alertContactCtl.GetAlertContact)
		group.POST("/setAlertContact", alertContactCtl.InsertAlertContact)
		group.POST("/updateAlertContact", alertContactCtl.UpdateAlertContact)
		group.POST("/deleteAlertContact", alertContactCtl.DeleteAlertContact)
		group.GET("/certifyAlertContact", alertContactCtl.CertifyAlertContact)
	}
}

func alertContactGroupRouters() {
	alertContactGroupCtl := controllers.NewAlertContactGroupCtl(dao.NewAlertContactGroup(database.GetDb()))
	group := router.Group("/hawkeye/alertContactGroup/")
	{
		group.GET("/getAlertContactGroup", alertContactGroupCtl.GetAlertContactGroup)
		group.GET("/getAlertContact", alertContactGroupCtl.GetAlertGroupContact)
		group.POST("/setAlertContactGroup", alertContactGroupCtl.InsertAlertContactGroup)
		group.POST("/updateAlertContactGroup", alertContactGroupCtl.UpdateAlertContactGroup)
		group.POST("/deleteAlertContactGroup", alertContactGroupCtl.DeleteAlertContactGroup)
	}
}
