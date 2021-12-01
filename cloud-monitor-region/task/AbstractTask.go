package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/producer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
)

var pageSize = 10000
var pageIndex = 1

type InterfaceTask interface {
	GetInstanceList(tenantId string) ([]*models.AlarmInstance, error)
}
type ProjectInstanceInfo struct {
	InstanceName string
	InstanceId   string
}

func SyncUpdate(task InterfaceTask, productType string, update bool) {
	var index = 1
	var pageTotal = 1
	for index <= pageTotal {
		pageVo := dao.AlarmInstance.SelectTenantIdList(productType, index, pageSize)
		if pageVo.Total < 0 {
			break
		}
		pageTotal = pageVo.Total
		if pageTotal == 0 {
			return
		}
		tenantList := pageVo.Records.(*[]string)
		for _, tenantId := range *tenantList {
			dbInstanceList := dao.AlarmInstance.SelectInstanceList(tenantId, productType)
			instances, err := task.GetInstanceList(tenantId)
			if err != nil {
				continue
			}
			if len(instances) > 0 {
				if update {
					dao.AlarmInstance.UpdateBatchInstanceName(instances)
					producer.SendInstanceJobMsg(sysRocketMq.InstanceTopic, instances)
				}
			}
			DeleteNotExistsInstances(tenantId, dbInstanceList, instances)
		}
		index++
	}
}

func DeleteNotExistsInstances(tenantId string, dbInstanceList []*models.AlarmInstance, instanceInfoList []*models.AlarmInstance) {
	var deletedList []*models.AlarmInstance
	for _, oldInstance := range dbInstanceList {
		exist := false
		for _, newInstance := range instanceInfoList {
			if IsEqual(oldInstance, newInstance) {
				exist = true
				break
			}
		}
		if !exist {
			deletedList = append(deletedList, oldInstance)
		}
	}
	if len(deletedList) != 0 {
		dao.AlarmInstance.DeleteInstanceList(tenantId, deletedList)
		instance := &dtos.Instance{
			TenantId: tenantId,
			List:     deletedList,
		}
		producer.SendInstanceJobMsg(sysRocketMq.DeleteInstanceTopic, instance)
		service.PrometheusRule.GenerateUserPrometheusRule(tenantId)
	}
}

func IsEqual(A, B interface{}) bool {
	if _, ok := A.(*models.AlarmInstance); ok {
		if _, ok := B.(*models.AlarmInstance); ok {
			if A.(*models.AlarmInstance).InstanceID == B.(*models.AlarmInstance).InstanceID {
				return A.(*models.AlarmInstance).InstanceName == B.(*models.AlarmInstance).InstanceName
			} else {
				return false
			}
		}
		return false
	}
	return false
}
