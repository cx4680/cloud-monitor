package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	actuator2 "code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/actuator"
	controller2 "code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/controller"
	inner2 "code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/inner"
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
	alarmRuleRouters()
	instanceRouters()
	alertRecordRouters()
	actuatorMapping()
	configItemRouters()
	noticeRouters()
	innerMapping()
}

func monitorProductRouters() {
	monitorProductCtl := controller2.NewMonitorProductCtl(commonDao.MonitorProduct)
	group := router.Group("/hawkeye/monitorProduct/")
	{
		group.GET("/getAllMonitorProducts", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetAllMonitorProducts)
		group.GET("/getById", monitorProductCtl.GetById)
		group.PUT("/updateById", monitorProductCtl.UpdateById)
	}
}

func monitorItemRouters() {
	monitorItemCtl := controller2.NewMonitorItemCtl(commonDao.MonitorItem)
	group := router.Group("/hawkeye/monitorItem/")
	{
		group.GET("/getMonitorItemsById", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.GetMonitorItemsById)
	}
}

func alertContactRouters() {
	alertContactCtl := controller2.NewAlertContactCtl(service.AlertContactService{})
	group := router.Group("/hawkeye/alertContact/")
	{
		group.GET("/getAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.GetAlertContact)
		group.POST("/setAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.InsertAlertContact)
		group.POST("/updateAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.UpdateAlertContact)
		group.POST("/deleteAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.DeleteAlertContact)
		group.GET("/certifyAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CertifyAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.CertifyAlertContact)
	}
}

func alertContactGroupRouters() {
	alertContactGroupCtl := controller2.NewAlertContactGroupCtl(service.AlertContactGroupService{})
	group := router.Group("/hawkeye/alertContactGroup/")
	{
		group.GET("/getAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.GetAlertContactGroup)
		group.GET("/getAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.GetAlertGroupContact)
		group.POST("/setAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.InsertAlertContactGroup)
		group.POST("/updateAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.UpdateAlertContactGroup)
		group.POST("/deleteAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.DeleteAlertContactGroup)
	}
}

func alarmRuleRouters() {
	ruleCtl := controller2.NewAlarmRuleCtl()
	group := router.Group("/hawkeye/rule/")
	{
		group.POST("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRulePageList", ResourceType: "*", ResourceId: "*"}), ruleCtl.SelectRulePageList)
		group.POST("/detail", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRuleDetail", ResourceType: "*", ResourceId: "*"}), ruleCtl.GetDetail)
		group.POST("/create", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CreateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.CreateRule)
		group.POST("/update", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.UpdateRule)
		group.POST("/delete", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.DeleteRule)
		group.POST("/changeStatus", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "ChangeGetAlertRuleStatus", ResourceType: "*", ResourceId: "*"}), ruleCtl.ChangeRuleStatus)
	}
}

func instanceRouters() {
	ctl := controller2.NewInstanceCtl()
	group := router.Group("/hawkeye/instance/")
	{
		group.POST("/rulePage", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRulePageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
		group.POST("/unbind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UnbindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Unbind)
		group.POST("/bind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "BindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Bind)
		group.POST("/ruleList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRuleList", ResourceType: "*", ResourceId: "*"}), ctl.GetRuleList)
	}
}

func alertRecordRouters() {
	ctl := controller2.NewAlertRecordController()
	group := router.Group("/hawkeye/alertRecord/")
	{
		group.POST("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordPageList", ResourceType: "*", ResourceId: "*"}), ctl.GetPageList)
		group.GET("/detail", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordDetail", ResourceType: "*", ResourceId: "*"}), ctl.GetDetail)
		group.GET("/total", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordTotal", ResourceType: "*", ResourceId: "*"}), ctl.GetAlertRecordTotal)
		group.GET("/recordNumHistory", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordNumHistory", ResourceType: "*", ResourceId: "*"}), ctl.GetRecordNumHistory)
	}

}

func actuatorMapping() {
	group := router.Group("/actuator")
	{
		group.GET("/env", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator2.Env())
		})
		group.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator2.Info())
		})
		group.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator2.Health())
		})
		group.GET("/metrics", func(c *gin.Context) {
			c.JSON(http.StatusOK, actuator2.Metrics())
		})
	}
}

func configItemRouters() {
	ctl := controller2.NewConfigItemCtl()
	group := router.Group("/hawkeye/configItem/")
	{
		group.GET("/getStatisticalPeriodList", ctl.GetStatisticalPeriodList)
		group.GET("/getContinuousCycleList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetContinuousCycleList", ResourceType: "*", ResourceId: "*"}), ctl.GetContinuousCycleList)
		group.GET("/getStatisticalMethodsList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetStatisticalMethodsList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalMethodsList)
		group.GET("/getComparisonMethodList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetComparisonMethodList", ResourceType: "*", ResourceId: "*"}), ctl.GetComparisonMethodList)
		group.GET("/getOverviewItemList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetOverviewItemList", ResourceType: "*", ResourceId: "*"}), ctl.GetOverviewItemList)
		group.GET("/getRegionList", ctl.GetRegionList)
		group.GET("/getMonitorRange", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorRangeList", ResourceType: "*", ResourceId: "*"}), ctl.GetMonitorRange)
		group.GET("/getNoticeChannel", ctl.GetNoticeChannel)
	}
}

func noticeRouters() {
	ctl := controller2.NewNoticeCtl(commonService.MessageService{})
	group := router.Group("/rest/notice/")
	{
		group.GET("/getUsage", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNoticeUsage", ResourceType: "*", ResourceId: "*"}), ctl.GetUsage)
	}
}

func innerMapping() {
	configItemController := inner2.NewConfigItemController()
	monitorItemController := inner2.NewMonitorItemController()
	ruleCtl := controller2.NewAlarmRuleCtl()
	innerRuleCtl := inner2.NewAlarmRuleCtl()
	group := router.Group("/hawkeye/inner/")
	{

		group.GET("configItem/getItemList", configItemController.GetItemListById)
		group.GET("monitorItem/getMonitorItemList", monitorItemController.GetMonitorItemsById)

		ruleGroup := group.Group("rule/")
		ruleGroup.POST("create", innerRuleCtl.CreateRule)
		ruleGroup.POST("update", innerRuleCtl.UpdateRule)
		ruleGroup.POST("delete", ruleCtl.DeleteRule)
		ruleGroup.POST("changeStatus", ruleCtl.ChangeRuleStatus)
	}
}
