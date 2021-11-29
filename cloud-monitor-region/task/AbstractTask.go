package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/producer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
)

var pageSize = 10000
var pageIndex = 1

type InterfaceTask interface {
	GetInstanceList(tenantId string) (interface{}, error)
}
type ProjectInstanceInfo struct {
	InstanceName string
	InstanceId   string
}

func SyncUpdate(task InterfaceTask, productType string, update bool) {
	var index = 1
	for {
		tenantIdList := dao.AlarmInstance.SelectTenantIdList(productType, index, pageSize)
		if len(tenantIdList) == 0 {
			break
		}
		for _, tenantId := range tenantIdList {
			dbInstanceList := dao.AlarmInstance.SelectInstanceList(tenantId, productType)
			instanceInfoList, err := task.GetInstanceList(tenantId)
			if err != nil {
				continue
			}
			infos := instanceInfoList.([]*ProjectInstanceInfo)
			var instances []*models.AlarmInstance
			if len(infos) > 0 {
				for _, instanceInfo := range infos {
					var alarmInstance = &models.AlarmInstance{
						InstanceID:   instanceInfo.InstanceId,
						InstanceName: instanceInfo.InstanceName,
					}
					instances = append(instances, alarmInstance)
				}
				if update {
					dao.AlarmInstance.UpdateBatchInstanceName(instances)
					producer.SendInstanceJobMsg(instances)
				}
			}
			DeleteNotExistsInstances(tenantId, dbInstanceList, instances)
		}
		index++
	}
}

func DeleteNotExistsInstances(tenantId string, dbInstanceList []*models.AlarmInstance, instanceInfoList []*models.AlarmInstance) {
	for i := len(dbInstanceList) - 1; i >= 0; i-- {
		v := dbInstanceList[i]
		for _, vv := range instanceInfoList {
			if v == vv {
				dbInstanceList = append(dbInstanceList[:i], dbInstanceList[i+1:]...)
			}
		}
	}
	if len(dbInstanceList) != 0 {
		dao.AlarmInstance.DeleteInstanceList(tenantId, dbInstanceList)
		service.PrometheusRule.GenerateUserPrometheusRule("", "", tenantId)
	}
}
