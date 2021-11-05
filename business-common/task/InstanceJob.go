package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"github.com/robfig/cron"
	"log"
)

var pageSize = "100"

var productType = "云服务器ECS"

func CronInstanceJob() {
	c := cron.New()
	err := c.AddFunc("0 0 0/1 * * ?", func() {
		instanceJob()
	})
	if err != nil {
		log.Println("clearAlertRecordJob error", err)
	}
	c.Start()
	defer c.Stop()
	select {}
}

func instanceJob() {
	log.Println("instanceJob start")
	var model = models.AlarmInstance{
		AlarmRuleID:  "1",
		InstanceID:   "1",
		InstanceName: "1",
		TenantID:     "1",
	}
	var model2 = models.AlarmInstance{
		AlarmRuleID:  "2",
		InstanceID:   "2",
		InstanceName: "2",
		TenantID:     "2",
	}
	var models = []models.AlarmInstance{model, model2}
	dao.NewAlarmInstanceDao(database.GetDb()).UpdateBatchInstanceName(models)
}
