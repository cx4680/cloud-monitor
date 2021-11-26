package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/controllers"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
)

func loadRouters() {
	monitorProductRouters()
	monitorItemRouters()
	alertContactRouters()
	alertContactGroupRouters()
	alarmRule()
	instance()
}

func monitorProductRouters() {
	monitorProductCtl := controllers.NewMonitorProductCtl(dao.MonitorProduct)
	group := router.Group("/hawkeye/monitorProduct/")
	{
		group.GET("/getAllMonitorProducts", monitorProductCtl.GetAllMonitorProducts)
		group.GET("/getById", monitorProductCtl.GetById)
		group.PUT("/updateById", monitorProductCtl.UpdateById)
	}
}

func monitorItemRouters() {
	monitorItemCtl := controllers.NewMonitorItemCtl(dao.MonitorItem)
	group := router.Group("/hawkeye/monitorItem/")
	{
		group.GET("/getMonitorItemsById", monitorItemCtl.GetMonitorItemsById)
	}
}

func alertContactRouters() {
	alertContactCtl := controllers.NewAlertContactCtl(service.AlertContactService{})
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
	alertContactGroupCtl := controllers.NewAlertContactGroupCtl(dao.AlertContactGroup)
	group := router.Group("/hawkeye/alertContactGroup/")
	{
		group.GET("/getAlertContactGroup", alertContactGroupCtl.GetAlertContactGroup)
		group.GET("/getAlertContact", alertContactGroupCtl.GetAlertGroupContact)
		group.POST("/setAlertContactGroup", alertContactGroupCtl.InsertAlertContactGroup)
		group.POST("/updateAlertContactGroup", alertContactGroupCtl.UpdateAlertContactGroup)
		group.POST("/deleteAlertContactGroup", alertContactGroupCtl.DeleteAlertContactGroup)
	}
}

func alarmRule() {
	ruleCtl := controllers.NewAlarmRuleCtl(commonDao.AlarmRule)
	group := router.Group("/hawkeye/rule/")
	{
		group.POST("/page", ruleCtl.SelectRulePageList)
		group.POST("/detail", ruleCtl.GetDetail)
		group.POST("/create", ruleCtl.CreateRule)
		group.POST("/update", ruleCtl.UpdateRule)
		group.POST("/delete", ruleCtl.DeleteRule)
		group.POST("/changeStatus", ruleCtl.ChangeRuleStatus)
	}
}

func instance() {
	ctl := controllers.NewInstanceCtl(commonDao.Instance)
	group := router.Group("/hawkeye/instance/")
	{
		group.POST("/rulePage", ctl.Page)
		group.POST("/unbind", ctl.Unbind)
		group.POST("/bind", ctl.Bind)
		group.POST("/ruleList", ctl.GetRuleList)
	}
}
