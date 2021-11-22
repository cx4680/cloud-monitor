package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"github.com/robfig/cron"
	"log"
)

var pageSize = 100

var productType = "云服务器ECS"

func CronInstanceJob() {
	c := cron.New()
	err := c.AddFunc("0 0 0/1 * * ?", instanceJob)
	if err != nil {
		log.Println("clearAlertRecordJob error", err)
	}
	c.Start()
	defer c.Stop()
	select {}
}

func instanceJob() {
	log.Println("instanceJob start")
	syncUpdate()
	log.Println("instanceJob end")
}

func syncUpdate() {
	var index = 1
	alarmInstanceDao := dao.AlarmInstance
	for {
		tenantIdList := alarmInstanceDao.SelectTenantIdList(productType, index, pageSize)
		if len(tenantIdList) == 0 {
			break
		}
		for _, tenantId := range tenantIdList {
			dbInstanceList := alarmInstanceDao.SelectInstanceList(tenantId, productType)
			var pageForm = forms.EcsQueryPageForm{
				TenantId: tenantId,
				Current:  1,
				PageSize: 10000,
			}
			pageVO := service.EcsPageList(pageForm)
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
					alarmInstanceDao.UpdateBatchInstanceName(instances)
					//TODO mq
				}
				DeleteNotExistsInstances(tenantId, dbInstanceList, instances)
			}
		}
		index++
	}
}
