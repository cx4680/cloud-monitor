package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
)

type InstanceDao struct {
}

var Instance = new(InstanceDao)

func (dao *InstanceDao) GetInstanceNum(tenantId string, regionCode string) int {
	var result int
	global.DB.Raw("select COUNT(DISTINCT arr.resource_id) num FROM t_alarm_rule_resource_rel arr  JOIN t_alarm_rule ar ON ar.biz_id = arr.alarm_rule_id  join t_resource res on res.instance_id=arr.resource_id  WHERE ar.product_name = '云服务器ECS'  AND ar.tenant_id = ?  AND res.region_code=? and ar.enabled=1", tenantId, regionCode).Scan(&result)
	return result
}
