package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
)

type InstanceDao struct {
}

var Instance = new(InstanceDao)

func (dao *InstanceDao) GetInstanceNum(tenantId string) int {
	var result int
	global.DB.Raw(" SELECT COUNT(DISTINCT arr.resource_id) num FROM t_alarm_rule_resource_rel arr  JOIN t_alarm_rule ar ON ar.id = arr.alarm_rule_id  WHERE product_type = '云服务器ECS' AND ar.tenant_id = ?", tenantId).Scan(&result)
	return result
}
