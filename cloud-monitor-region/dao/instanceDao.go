package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
)

type InstanceDao struct {
}

var Instance = new(InstanceDao)

func (dao *InstanceDao) GetInstanceNum(tenantId string) int {
	var result int
	global.DB.Raw(" SELECT count(DISTINCT arr.resource_id) num from t_alarm_rule_resource_rel arr  join t_alarm_rule ar on ar.id = arr.alarm_rule_id  where ar.tenant_id=?", tenantId).Scan(&result)
	return result
}
