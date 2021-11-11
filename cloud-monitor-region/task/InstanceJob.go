package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
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

func syncUpdate() {
	var pages = 1
	var index = 1
	instanceDao := dao.NewAlarmInstanceDao(database.GetDb())
	for index <= pages {
		tenantIdList := instanceDao.SelectTenantIdList(productType, index, pages)
		for _, tenantId := range tenantIdList {
			//instanceList := instanceDao.SelectInstanceList(tenantId, productType)
			var pageForm = forms.EcsQueryPageForm{
				TenantId: tenantId,
				Current:  1,
				PageSize: 10000,
			}
			pageVO := service.EcsPageList(tenantId, pageForm)
			if pageVO.Total > 0 {
				instanceInfoList := pageVO.Rows
				var instances []models.AlarmInstance
				if len(instanceInfoList) > 0 {
					for _, instanceInfo := range instanceInfoList {
						var alarmInstance = models.AlarmInstance{
							InstanceID:   instanceInfo.InstanceId,
							InstanceName: instanceInfo.InstanceName,
						}
						instances = append(instances, alarmInstance)
					}
					instanceDao.UpdateBatchInstanceName(instances)
					//TODO mq
				}
			}
		}
	}
}
