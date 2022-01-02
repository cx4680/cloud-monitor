package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/actuator"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/controller"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/inner"
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
	monitorProductCtl := controller.NewMonitorProductCtl(commonDao.MonitorProduct)
	group := router.Group("/hawkeye/monitorProduct/")
	{
		group.GET("/getAllMonitorProducts", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetAllMonitorProducts)
	}
}

func monitorItemRouters() {
	monitorItemCtl := controller.NewMonitorItemCtl(commonDao.MonitorItem)
	group := router.Group("/hawkeye/monitorItem/")
	{
		group.GET("/getMonitorItemsById", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.GetMonitorItemsById)
	}
}

func alertContactRouters() {
	alertContactCtl := controller.NewAlertContactCtl(service.AlertContactService{})
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
	alertContactGroupCtl := controller.NewAlertContactGroupCtl(service.AlertContactGroupService{})
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
	ruleCtl := controller.NewAlarmRuleCtl()
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
	ctl := controller.NewInstanceCtl()
	group := router.Group("/hawkeye/instance/")
	{
		group.POST("/rulePage", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRulePageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
		group.POST("/unbind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UnbindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Unbind)
		group.POST("/bind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "BindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Bind)
		group.POST("/ruleList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRuleList", ResourceType: "*", ResourceId: "*"}), ctl.GetRuleList)
	}
}

func alertRecordRouters() {
	ctl := controller.NewAlertRecordController()
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

func configItemRouters() {
	ctl := controller.NewConfigItemCtl()
	group := router.Group("/hawkeye/configItem/")
	{
		group.GET("/getStatisticalPeriodList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetStatisticalPeriodList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalPeriodList)
		group.GET("/getContinuousCycleList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetContinuousCycleList", ResourceType: "*", ResourceId: "*"}), ctl.GetContinuousCycleList)
		group.GET("/getStatisticalMethodsList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetStatisticalMethodsList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalMethodsList)
		group.GET("/getComparisonMethodList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetComparisonMethodList", ResourceType: "*", ResourceId: "*"}), ctl.GetComparisonMethodList)
		group.GET("/getOverviewItemList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetOverviewItemList", ResourceType: "*", ResourceId: "*"}), ctl.GetOverviewItemList)
		group.GET("/getRegionList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetRegionList", ResourceType: "*", ResourceId: "*"}), ctl.GetRegionList)
		group.GET("/getMonitorRange", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorRangeList", ResourceType: "*", ResourceId: "*"}), ctl.GetMonitorRange)
		group.GET("/getNoticeChannel", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNoticeChannelList", ResourceType: "*", ResourceId: "*"}), ctl.GetNoticeChannel)
	}
}

func noticeRouters() {
	ctl := controller.NewNoticeCtl(commonService.MessageService{})
	group := router.Group("/rest/notice/")
	{
		group.GET("/getUsage", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNoticeUsage", ResourceType: "*", ResourceId: "*"}), ctl.GetUsage)
	}
}

func innerMapping() {
	configItemController := inner.NewConfigItemController()
	monitorItemController := inner.NewMonitorItemController()
	ruleCtl := controller.NewAlarmRuleCtl()
	innerRuleCtl := inner.NewAlarmRuleCtl()
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
