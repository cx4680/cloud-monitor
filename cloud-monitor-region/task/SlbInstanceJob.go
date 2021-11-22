package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"github.com/robfig/cron"
	"log"
)

var slbProductType = "负载均衡SLB"

func CronSlbInstanceJob() {
	c := cron.New()
	err := c.AddFunc("0 0 0/1 * * ?", slbInstanceJob)
	if err != nil {
		log.Println("clearAlertRecordJob error", err)
	}
	c.Start()
	defer c.Stop()
	select {}
}

func slbInstanceJob() {
	log.Println("slbInstanceJob start")
	slbSyncUpdate()
	log.Println("slbInstanceJob end")
}

func slbSyncUpdate() {
	var index = 1
	alarmInstanceDao := dao.AlarmInstance
	for {
		tenantIdList := alarmInstanceDao.SelectTenantIdList(slbProductType, index, pageSize)
		if len(tenantIdList) == 0 {
			break
		}
		for _, tenantId := range tenantIdList {
			dbInstanceList := alarmInstanceDao.SelectInstanceList(tenantId, slbProductType)
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
					mq.SendMsg(config.GetRocketmqConfig().InstanceTopic, tools.ToString(instances))
				}
				DeleteNotExistsInstances(tenantId, dbInstanceList, instances)
			}
		}
		index++
	}
}
