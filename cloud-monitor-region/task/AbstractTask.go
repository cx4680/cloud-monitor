package task

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
)

func DeleteNotExistsInstances(tenantId string, dbInstanceList []models.AlarmInstance, instanceInfoList []models.AlarmInstance) {
	for i := len(dbInstanceList) - 1; i >= 0; i-- {
		v := dbInstanceList[i]
		for _, vv := range instanceInfoList {
			if v == vv {
				dbInstanceList = append(dbInstanceList[:i], dbInstanceList[i+1:]...)
			}
		}
	}
	if len(dbInstanceList) != 0 {
		commonDao.AlarmInstance.DeleteInstanceList(tenantId, dbInstanceList)
		dao.PrometheusRule.GenerateUserPrometheusRule("", "", tenantId)
	}
}
