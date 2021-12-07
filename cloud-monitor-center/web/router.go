package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/actuator"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/controllers"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func loadRouters() {
	monitorProductRouters()
	monitorItemRouters()
	alertContactRouters()
	alertContactGroupRouters()
	alarmRule()
	instance()
	alertRecord()

	actuatorMapping()
}

func monitorProductRouters() {
	monitorProductCtl := controllers.NewMonitorProductCtl(commonDao.MonitorProduct)
	group := router.Group("/hawkeye/monitorProduct/")
	{
		group.GET("/getAllMonitorProducts", global.IamAuthIdentify(&models.Identity{Product: "ECS", Action: "GetInstanceList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetAllMonitorProducts)
		group.GET("/getById", monitorProductCtl.GetById)
		group.PUT("/updateById", monitorProductCtl.UpdateById)
	}
}

func monitorItemRouters() {
	monitorItemCtl := controllers.NewMonitorItemCtl(commonDao.MonitorItem)
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
	alertContactGroupCtl := controllers.NewAlertContactGroupCtl(service.AlertContactGroupService{})
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

func alertRecord() {
	ctl := controllers.NewAlertRecordController()
	group := router.Group("/hawkeye/alertRecord/")
	{
		group.POST("/page", ctl.GetPageList)
		group.GET("/detail", ctl.GetDetail)
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
