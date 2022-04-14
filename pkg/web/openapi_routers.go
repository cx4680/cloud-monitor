package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/api/openapi/v1.0"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/iam"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/logs"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
)

func loadOpenApiV1Routers() {
	group := Router.Group("/v1.0/")
	monitorProductOpenApiV1Routers(group)
	monitorItemOpenApiV1Routers(group)
	contactOpenApiV1Routers(group)
	contactGroupOpenApiV1Routers(group)
	instanceOpenApiRouters(group)
	ruleOpenApiRouters(group)
	alarmHistoryOpiRouters(group)

	MonitorReportOpenApiV1Routers(group)
	ResourceOpenApiV1Routers(group)
}

func monitorProductOpenApiV1Routers(group *gin.RouterGroup) {
	monitorProductCtl := v1_0.NewMonitorProductCtl(service.MonitorProductService{})
	group.GET("products", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetMonitorProduct)
}

func monitorItemOpenApiV1Routers(group *gin.RouterGroup) {
	monitorItemCtl := v1_0.NewMonitorItemCtl(service.MonitorItemService{})
	group.GET("products/:ProductAbbreviation/metrics", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.GetMonitorItemsByProductAbbr)
}

func instanceOpenApiRouters(group *gin.RouterGroup) {
	ctl := v1_0.NewInstanceCtl()
	group.GET("resources/:ResourceId/rules", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRulePageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
	group.DELETE("resources/:ResourceId/rules", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UnbindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Unbind)
	group.PUT("resources/:ResourceId/rules", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "BindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Bind)
}

func contactOpenApiV1Routers(group *gin.RouterGroup) {
	contactCtl := v1_0.NewContactCtl()
	group.GET("contacts", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.SelectContactPage)
	group.POST("contacts", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.CreateContact)
	group.PUT("contacts/:ContactId", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.UpdateContact)
	group.DELETE("contacts/:ContactId", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.DeleteContact)
	group.PUT("contacts/activate/:ActiveCode", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CertifyAlertContact", ResourceType: "*", ResourceId: "*"}), contactCtl.ActivateContact)
}

func contactGroupOpenApiV1Routers(group *gin.RouterGroup) {
	contactGroupCtl := v1_0.NewContactGroupCtl()
	group.GET("groups", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.SelectContactGroupPage)
	group.GET("groups/:GroupId/contacts", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertContact", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.SelectContactPageByGroupId)
	group.POST("groups", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "SetAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.CreateContactGroup)
	group.PUT("groups/:GroupId", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.UpdateContactGroup)
	group.DELETE("groups/:GroupId", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertContactGroup", ResourceType: "*", ResourceId: "*"}), contactGroupCtl.DeleteContactGroup)
}

func ruleOpenApiRouters(group *gin.RouterGroup) {
	ruleCtl := v1_0.NewAlarmRuleCtl()
	{
		group.GET("rules", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRulePageList", ResourceType: "*", ResourceId: "*"}), ruleCtl.SelectRulePageList)
		group.GET("rules/:RuleId", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRuleDetail", ResourceType: "*", ResourceId: "*"}), ruleCtl.GetDetail)
		group.POST("rules", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "CreateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.CreateRule)
		group.PUT("rules/:RuleId", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UpdateAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.UpdateRule)
		group.DELETE("rules/:RuleId", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "DeleteAlertRule", ResourceType: "*", ResourceId: "*"}), ruleCtl.DeleteRule)
		group.PUT("rules/:RuleId/status", logs.GinTrailzap(false, Write), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "ChangeGetAlertRuleStatus", ResourceType: "*", ResourceId: "*"}), ruleCtl.ChangeRuleStatus)
	}
}

func alarmHistoryOpiRouters(group *gin.RouterGroup) {
	ctl := v1_0.NewAlarmRecordController()
	{
		group.GET("alarms", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordPageList", ResourceType: "*", ResourceId: "*"}), ctl.GetPageList)
		group.GET("alarms/:AlarmBizId/contacts", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetContactInfoList", ResourceType: "*", ResourceId: "*"}), ctl.GetAlarmContactInfo)
	}
}

func MonitorReportOpenApiV1Routers(group *gin.RouterGroup) {
	monitorReportFormCtl := v1_0.NewMonitorReportFormController()
	group.GET("resources/:ResourceId/metrics/:MetricCode/datas", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportRangeData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetMonitorDatas)
	group.GET("resources/:ResourceId/metrics/:MetricCode/data", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportData", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetMonitorData)
	group.GET("metrics/:MetricCode/:N/resources", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorReportTop", ResourceType: "*", ResourceId: "*"}), monitorReportFormCtl.GetMonitorDataTop)
}

func ResourceOpenApiV1Routers(group *gin.RouterGroup) {
	resourceCtl := v1_0.NewResourceController()
	group.GET(":ProductAbbreviation/resources", logs.GinTrailzap(false, Read), iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetResourceList", ResourceType: "*", ResourceId: "*"}), resourceCtl.GetResourceList)
}
