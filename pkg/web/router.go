package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/api/actuator"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/api/controller"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/api/inner"
	iam2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/iam"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/logs"
	service2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/task"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const Read = "Read"
const Write = "Write"

const pathPrefix = "/hawkeye/"

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

	instance()
	MonitorReportForm()
	innerCtl()
	remote()
}

func monitorProductRouters() {
	monitorProductCtl := controller.NewMonitorProductCtl(service.MonitorProductService{})
	group := Router.Group(pathPrefix + "monitorProduct/")
	{
		group.GET("/getAllMonitorProducts", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetMonitorProduct)
		group.GET("/getMonitorProduct", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetAllMonitorProduct)
		group.POST("/changeStatus", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.ChangeStatus)
	}
}

func monitorItemRouters() {
	monitorItemCtl := controller.NewMonitorItemCtl(service.MonitorItemService{})
	group := Router.Group(pathPrefix + "monitorItem/")
	{
		group.GET("/getMonitorItemsById", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.GetMonitorItemsById)
		group.POST("/changeDisplay", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.ChangeDisplay)
	}
}

func contactRouters() {
	contactCtl := controller.NewContactCtl()
	group := Router.Group(pathPrefix + "contact/")
	{
		group.GET("/getContact", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.GetContact)
		group.POST("/addContact", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "SetAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.CreateContact)
		group.POST("/updateContact", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "UpdateAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.UpdateContact)
		group.POST("/deleteContact", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "DeleteAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.DeleteContact)
		group.GET("/activateContact", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "CertifyAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.ActivateContact)
	}
}

func contactGroupRouters() {
	contactGroupCtl := controller.NewContactGroupCtl()
	group := Router.Group(pathPrefix + "contactGroup/")
	{
		group.GET("/getContactGroup", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.GetContactGroup)
		group.GET("/getContact", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.GetGroupContact)
		group.POST("/addContactGroup", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "SetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.CreateContactGroup)
		group.POST("/updateContactGroup", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "UpdateAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.UpdateContactGroup)
		group.POST("/deleteContactGroup", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "DeleteAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.DeleteContactGroup)
	}
}

func alarmRuleRouters() {
	ruleCtl := controller.NewAlarmRuleCtl()
	group := Router.Group(pathPrefix + "rule/")
	{
		group.POST("/page", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlertRulePageList", ResourceType: "*", ResourceId: "*"}), ruleCtl.SelectRulePageList)
		group.POST("/detail", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlertRuleDetail", ResourceType: "*", ResourceId: "*"}), ruleCtl.GetDetail)
		group.POST("/create", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "CreateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.CreateRule)
		group.POST("/update", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "UpdateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.UpdateRule)
		group.POST("/delete", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "DeleteAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.DeleteRule)
		group.POST("/changeStatus", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "ChangeGetAlertRuleStatus", ResourceType: "*", ResourceId: "*"}), ruleCtl.ChangeRuleStatus)
	}
}

func instanceRouters() {
	ctl := controller.NewInstanceCtl()
	group := Router.Group(pathPrefix + "instance/")
	{
		group.POST("/rulePage", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetInstanceRulePageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
		group.POST("/unbind", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "UnbindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Unbind)
		group.POST("/bind", logs.GinTrailzap(false, Write), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "BindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Bind)
		group.POST("/ruleList", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetInstanceRuleList", ResourceType: "*", ResourceId: "*"}), ctl.GetRuleList)
	}
}

func alertRecordRouters() {
	ctl := controller.NewAlarmRecordController()
	group := Router.Group(pathPrefix + "alarmRecord/")
	{
		group.POST("/page", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlertRecordPageList", ResourceType: "*", ResourceId: "*"}), ctl.GetPageList)
		group.GET("/contactInfos", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetContactInfoList", ResourceType: "*", ResourceId: "*"}), ctl.GetAlarmContactInfo)
		group.GET("/total", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlarmRecordTotal", ResourceType: "*", ResourceId: "*"}), ctl.GetAlarmRecordTotal)
		group.GET("/recordNumHistory", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetAlertRecordNumHistory", ResourceType: "*", ResourceId: "*"}), ctl.GetRecordNumHistory)
	}

}

func actuatorMapping() {
	group := Router.Group("/actuator")
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
	group := Router.Group(pathPrefix + "configItem/")
	{
		group.GET("/getStatisticalPeriodList", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetStatisticalPeriodList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalPeriodList)
		group.GET("/getContinuousCycleList", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetContinuousCycleList", ResourceType: "*", ResourceId: "*"}), ctl.GetContinuousCycleList)
		group.GET("/getStatisticalMethodsList", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetStatisticalMethodsList", ResourceType: "*", ResourceId: "*"}), ctl.GetStatisticalMethodsList)
		group.GET("/getComparisonMethodList", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetComparisonMethodList", ResourceType: "*", ResourceId: "*"}), ctl.GetComparisonMethodList)
		group.GET("/getOverviewItemList", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetOverviewItemList", ResourceType: "*", ResourceId: "*"}), ctl.GetOverviewItemList)
		group.GET("/getRegionList", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetRegionList", ResourceType: "*", ResourceId: "*"}), ctl.GetRegionList)
		group.GET("/getMonitorRange", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetMonitorRangeList", ResourceType: "*", ResourceId: "*"}), ctl.GetMonitorRange)
		group.GET("/getNoticeChannel", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetNoticeChannelList", ResourceType: "*", ResourceId: "*"}), ctl.GetNoticeChannel)
	}
}

func noticeRouters() {
	ctl := controller.NewNoticeCtl(service2.MessageService{})
	group := Router.Group(pathPrefix + "notice")
	{
		group.GET("/getUsage", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetNoticeUsage", ResourceType: "*", ResourceId: "*"}), ctl.GetUsage)
	}
}

func innerMapping() {
	configItemController := inner.NewConfigItemController()
	monitorItemController := inner.NewMonitorItemController()
	ruleCtl := controller.NewAlarmRuleCtl()
	innerRuleCtl := inner.NewAlarmRuleCtl()
	noticeCtl := controller.NewNoticeCtl(service2.MessageService{})
	group := Router.Group(pathPrefix + "inner/")
	{
		group.GET("configItem/getItemList", configItemController.GetItemListById)
		group.GET("monitorItem/getMonitorItemList", monitorItemController.GetMonitorItemsById)
		group.GET("notice/getUsage", noticeCtl.GetCenterUsage)
		group.POST("notice/changeNoticeChannel", noticeCtl.ChangeNoticeChannel)

		ruleGroup := group.Group("rule/")
		ruleGroup.POST("create", innerRuleCtl.CreateRule)
		ruleGroup.POST("update", innerRuleCtl.UpdateRule)
		ruleGroup.POST("delete", ruleCtl.DeleteRule)
		ruleGroup.POST("changeStatus", ruleCtl.ChangeRuleStatus)
	}
}

func MonitorReportForm() {
	monitorReportFormCtl := controller.NewMonitorReportFormController(service.NewMonitorReportFormService())
	group := Router.Group(pathPrefix + "MonitorReportForm/")
	{
		group.GET("/getData", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetMonitorReportData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetData)
		group.GET("/getAxisData", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetMonitorReportRangeData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetAxisData)
		group.GET("/getTop", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetMonitorReportTop", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetTop)
	}
}

func instance() {
	instanceCtl := controller.NewInstanceRegionCtl(dao.Instance)
	group := Router.Group(pathPrefix + "instance/")
	{
		group.GET("/page", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetInstancePageList", ResourceType: "*", ResourceId: "*"}), instanceCtl.GetPage)
		group.GET("/getInstanceNum", logs.GinTrailzap(false, Read), iam2.AuthIdentify(&models.Identity{Product: iam2.ProductMonitor, Action: "GetInstanceNum", ResourceType: "*", ResourceId: "*"}), instanceCtl.GetInstanceNumByRegion)
	}
}

func innerCtl() {
	addService := service.NewAlarmRecordAddService(service.NewAlarmRecordService(), service2.NewAlarmHandlerService(), service2.NewTenantService())
	ctl := inner.NewAlertRecordCtl(addService)
	group := Router.Group("/inner/")
	{
		group.POST("/alarmRecord/insert", ctl.AddAlarmRecord)
	}
}

func remote() {
	Router.GET("/inner/remote/:productType", func(context *gin.Context) {
		productType := context.Param("productType")
		task.Run(productType)
	})
}
