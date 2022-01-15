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
	contactRouters()
	contactGroupRouters()
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
		group.GET("/getAllMonitorProducts", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetAllMonitorProducts)
	}
}

func monitorItemRouters() {
	monitorItemCtl := controller.NewMonitorItemCtl(commonDao.MonitorItem)
	group := router.Group("/hawkeye/monitorItem/")
	{
		group.GET("/getMonitorItemsById", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.GetMonitorItemsById)
	}
}

func contactRouters() {
	contactCtl := controller.NewContactCtl(service.ContactService{})
	group := router.Group("/hawkeye/contact/")
	{
		group.GET("/getContact", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.GetContact)
		group.POST("/addContact", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.AddContact)
		group.POST("/updateContact", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.UpdateContact)
		group.POST("/deleteContact", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.DeleteContact)
		group.GET("/activateContact", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CertifyAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.ActivateContact)
	}
}

func contactGroupRouters() {
	contactGroupCtl := controller.NewContactGroupCtl(service.ContactGroupService{})
	group := router.Group("/hawkeye/contactGroup/")
	{
		group.GET("/getContactGroup", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.GetContactGroup)
		group.GET("/getContact", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.GetGroupContact)
		group.POST("/addContactGroup", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.AddContactGroup)
		group.POST("/updateContactGroup", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.UpdateContactGroup)
		group.POST("/deleteContactGroup", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.DeleteContactGroup)
	}
}

func alarmRuleRouters() {
	ruleCtl := controller.NewAlarmRuleCtl()
	group := router.Group("/hawkeye/rule/")
	{
		group.POST("/page", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRulePageList", ResourceType: "*", ResourceId: "*"}), ruleCtl.SelectRulePageList)
		group.POST("/detail", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRuleDetail", ResourceType: "*", ResourceId: "*"}), ruleCtl.GetDetail)
		group.POST("/create", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CreateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.CreateRule)
		group.POST("/update", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.UpdateRule)
		group.POST("/delete", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.DeleteRule)
		group.POST("/changeStatus", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "ChangeGetAlertRuleStatus", ResourceType: "*", ResourceId: "*"}), ruleCtl.ChangeRuleStatus)
	}
}

func instanceRouters() {
	ctl := controller.NewInstanceCtl()
	group := router.Group("/hawkeye/instance/")
	{
		group.POST("/rulePage", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRulePageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
		group.POST("/unbind", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UnbindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Unbind)
		group.POST("/bind", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "BindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Bind)
		group.POST("/ruleList", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRuleList", ResourceType: "*", ResourceId: "*"}), ctl.GetRuleList)
	}
}

func alertRecordRouters() {
	ctl := controller.NewAlarmRecordController()
	group := router.Group("/hawkeye/alarmRecord/")
	{
		group.POST("/page", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordPageList", ResourceType: "*", ResourceId: "*"}), ctl.GetPageList)
		group.GET("/contactInfos", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetContactInfoList", ResourceType: "*", ResourceId: "*"}), ctl.GetAlarmContactInfo)
		group.GET("/total", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlarmRecordTotal", ResourceType: "*", ResourceId: "*"}), ctl.GetAlarmRecordTotal)
		group.GET("/recordNumHistory", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordNumHistory", ResourceType: "*", ResourceId: "*"}), ctl.GetRecordNumHistory)
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
		group.GET("/getStatisticalPeriodList", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetStatisticalPeriodList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalPeriodList)
		group.GET("/getContinuousCycleList", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetContinuousCycleList", ResourceType: "*", ResourceId: "*"}), ctl.GetContinuousCycleList)
		group.GET("/getStatisticalMethodsList", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetStatisticalMethodsList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalMethodsList)
		group.GET("/getComparisonMethodList", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetComparisonMethodList", ResourceType: "*", ResourceId: "*"}), ctl.GetComparisonMethodList)
		group.GET("/getOverviewItemList", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetOverviewItemList", ResourceType: "*", ResourceId: "*"}), ctl.GetOverviewItemList)
		group.GET("/getRegionList", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetRegionList", ResourceType: "*", ResourceId: "*"}), ctl.GetRegionList)
		group.GET("/getMonitorRange", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorRangeList", ResourceType: "*", ResourceId: "*"}), ctl.GetMonitorRange)
		group.GET("/getNoticeChannel", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNoticeChannelList", ResourceType: "*", ResourceId: "*"}), ctl.GetNoticeChannel)
	}
}

func noticeRouters() {
	ctl := controller.NewNoticeCtl(commonService.MessageService{})
	group := router.Group("/rest/notice/")
	{
		group.GET("/getUsage", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetNoticeUsage", ResourceType: "*", ResourceId: "*"}), ctl.GetUsage)
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
