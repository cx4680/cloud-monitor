package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/logs"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/actuator"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/controller"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/inner"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const Read = "Read"
const Write = "Write"

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
		group.GET("/getAllMonitorProducts", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetAllMonitorProducts, logs.GinTrailzap(false, Read))
	}
}

func monitorItemRouters() {
	monitorItemCtl := controller.NewMonitorItemCtl(commonDao.MonitorItem)
	group := router.Group("/hawkeye/monitorItem/")
	{
		group.GET("/getMonitorItemsById", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.GetMonitorItemsById, logs.GinTrailzap(false, Read))
	}
}

func alertContactRouters() {
	alertContactCtl := controller.NewAlertContactCtl(service.AlertContactService{})
	group := router.Group("/hawkeye/alertContact/")
	{
		group.GET("/getAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.GetAlertContact, logs.GinTrailzap(false, Read))
		group.POST("/setAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.InsertAlertContact, logs.GinTrailzap(false, Write))
		group.POST("/updateAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.UpdateAlertContact, logs.GinTrailzap(false, Write))
		group.POST("/deleteAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.DeleteAlertContact, logs.GinTrailzap(false, Write))
		group.GET("/certifyAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CertifyAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactCtl.CertifyAlertContact, logs.GinTrailzap(false, Write))
	}
}

func alertContactGroupRouters() {
	alertContactGroupCtl := controller.NewAlertContactGroupCtl(service.AlertContactGroupService{})
	group := router.Group("/hawkeye/alertContactGroup/")
	{
		group.GET("/getAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.GetAlertContactGroup, logs.GinTrailzap(false, Read))
		group.GET("/getAlertContact", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.GetAlertGroupContact, logs.GinTrailzap(false, Read))
		group.POST("/setAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.InsertAlertContactGroup, logs.GinTrailzap(false, Write))
		group.POST("/updateAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.UpdateAlertContactGroup, logs.GinTrailzap(false, Write))
		group.POST("/deleteAlertContactGroup", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContactGroup", ResourceType: "*", ResourceId: "*"}), alertContactGroupCtl.DeleteAlertContactGroup, logs.GinTrailzap(false, Write))
	}
}

func alarmRuleRouters() {
	ruleCtl := controller.NewAlarmRuleCtl()
	group := router.Group("/hawkeye/rule/")
	{
		group.POST("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRulePageList", ResourceType: "*", ResourceId: "*"}), ruleCtl.SelectRulePageList, logs.GinTrailzap(false, Read))
		group.POST("/detail", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRuleDetail", ResourceType: "*", ResourceId: "*"}), ruleCtl.GetDetail, logs.GinTrailzap(false, Read))
		group.POST("/create", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CreateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.CreateRule, logs.GinTrailzap(false, Write))
		group.POST("/update", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.UpdateRule, logs.GinTrailzap(false, Write))
		group.POST("/delete", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.DeleteRule, logs.GinTrailzap(false, Write))
		group.POST("/changeStatus", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "ChangeGetAlertRuleStatus", ResourceType: "*", ResourceId: "*"}), ruleCtl.ChangeRuleStatus, logs.GinTrailzap(false, Write))
	}
}

func instanceRouters() {
	ctl := controller.NewInstanceCtl()
	group := router.Group("/hawkeye/instance/")
	{
		group.POST("/rulePage", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRulePageList", ResourceType: "*", ResourceId: "*"}), ctl.Page, logs.GinTrailzap(false, Read))
		group.POST("/unbind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UnbindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Unbind, logs.GinTrailzap(false, Write))
		group.POST("/bind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "BindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Bind, logs.GinTrailzap(false, Write))
		group.POST("/ruleList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRuleList", ResourceType: "*", ResourceId: "*"}), ctl.GetRuleList, logs.GinTrailzap(false, Read))
	}
}

func alertRecordRouters() {
	ctl := controller.NewAlertRecordController()
	group := router.Group("/hawkeye/alertRecord/")
	{
		group.POST("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordPageList", ResourceType: "*", ResourceId: "*"}), ctl.GetPageList, logs.GinTrailzap(false, Read))
		group.GET("/detail", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordDetail", ResourceType: "*", ResourceId: "*"}), ctl.GetDetail, logs.GinTrailzap(false, Read))
		group.GET("/contactInfos", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetContactInfoList", ResourceType: "*", ResourceId: "*"}), ctl.GetAlarmContactInfo, logs.GinTrailzap(false, Read))
		group.GET("/total", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordTotal", ResourceType: "*", ResourceId: "*"}), ctl.GetAlertRecordTotal, logs.GinTrailzap(false, Read))
		group.GET("/recordNumHistory", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordNumHistory", ResourceType: "*", ResourceId: "*"}), ctl.GetRecordNumHistory, logs.GinTrailzap(false, Read))
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
		group.GET("/getStatisticalPeriodList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetStatisticalPeriodList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalPeriodList, logs.GinTrailzap(false, Read))
		group.GET("/getContinuousCycleList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetContinuousCycleList", ResourceType: "*", ResourceId: "*"}), ctl.GetContinuousCycleList, logs.GinTrailzap(false, Read))
		group.GET("/getStatisticalMethodsList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetStatisticalMethodsList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalMethodsList, logs.GinTrailzap(false, Read))
		group.GET("/getComparisonMethodList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetComparisonMethodList", ResourceType: "*", ResourceId: "*"}), ctl.GetComparisonMethodList, logs.GinTrailzap(false, Read))
		group.GET("/getOverviewItemList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetOverviewItemList", ResourceType: "*", ResourceId: "*"}), ctl.GetOverviewItemList, logs.GinTrailzap(false, Read))
		group.GET("/getRegionList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetRegionList", ResourceType: "*", ResourceId: "*"}), ctl.GetRegionList, logs.GinTrailzap(false, Read))
		group.GET("/getMonitorRange", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorRangeList", ResourceType: "*", ResourceId: "*"}), ctl.GetMonitorRange, logs.GinTrailzap(false, Read))
		group.GET("/getNoticeChannel", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNoticeChannelList", ResourceType: "*", ResourceId: "*"}), ctl.GetNoticeChannel, logs.GinTrailzap(false, Read))
	}
}

func noticeRouters() {
	ctl := controller.NewNoticeCtl(commonService.MessageService{})
	group := router.Group("/rest/notice/")
	{
		group.GET("/getUsage", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNoticeUsage", ResourceType: "*", ResourceId: "*"}), ctl.GetUsage, logs.GinTrailzap(false, Read))
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
