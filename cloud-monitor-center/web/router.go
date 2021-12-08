package web

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/iam"
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
	alarmRuleRouters()
	instanceRouters()
	alertRecordRouters()
	actuatorMapping()
	configItemRouters()
}

func monitorProductRouters() {
	monitorProductCtl := controllers.NewMonitorProductCtl(commonDao.MonitorProduct)
	group := router.Group("/hawkeye/monitorProduct/")
	{
		group.GET("/getAllMonitorProducts", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAllMonitorProductsList", ResourceType: "*", ResourceId: "*"}), monitorProductCtl.GetAllMonitorProducts)
		group.GET("/getById", monitorProductCtl.GetById)
		group.PUT("/updateById", monitorProductCtl.UpdateById)
	}
}

func monitorItemRouters() {
	monitorItemCtl := controllers.NewMonitorItemCtl(commonDao.MonitorItem)
	group := router.Group("/hawkeye/monitorItem/")
	{
		group.GET("/getMonitorItemsById", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetMonitorItemsByIdList", ResourceType: "*", ResourceId: "*"}), monitorItemCtl.GetMonitorItemsById)
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

func alarmRuleRouters() {
	ruleCtl := controllers.NewAlarmRuleCtl(commonDao.AlarmRule)
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
	ctl := controllers.NewInstanceCtl(commonDao.Instance)
	group := router.Group("/hawkeye/instance/")
	{
		group.POST("/rulePage", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRulePageList", ResourceType: "*", ResourceId: "*"}), ctl.Page)
		group.POST("/unbind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "UnbindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Unbind)
		group.POST("/bind", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "BindInstanceRule", ResourceType: "*", ResourceId: "*"}), ctl.Bind)
		group.POST("/ruleList", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetInstanceRuleList", ResourceType: "*", ResourceId: "*"}), ctl.GetRuleList)
	}
}

func alertRecordRouters() {
	ctl := controllers.NewAlertRecordController()
	group := router.Group("/hawkeye/alertRecord/")
	{
		group.POST("/page", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordPageList", ResourceType: "*", ResourceId: "*"}), ctl.GetPageList)
		group.GET("/detail", iam.AuthIdentify(&models.Identity{Product: iam.ProductMonitor, Action: "GetAlertRecordDetail", ResourceType: "*", ResourceId: "*"}), ctl.GetDetail)
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
	ctl := controllers.NewConfigItemCtl()
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
