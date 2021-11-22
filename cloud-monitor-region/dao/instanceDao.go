package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
)

type InstanceDao struct {
}

var Instance = new(InstanceDao)

func (dao *InstanceDao) GetInstanceNum(tenantId string) int {
	var result int
	database.GetDb().Raw(" SELECT count(DISTINCT ai.instance_id) num from t_alarm_instance ai       join t_alarm_rule ar on ar.id=ai.alarm_rule_id       where ar.tenant_id=?", tenantId).Scan(&result)
	return result
}
