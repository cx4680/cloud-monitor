package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"fmt"
	"strings"
)

type AlarmInstanceDao struct {
}

var AlarmInstance = new(AlarmInstanceDao)

func (mpd *AlarmInstanceDao) SelectTenantIdList(productType string, pageCurrent int, pageSize int) *vo.PageVO {
	sql := "SELECT DISTINCT   t1.tenant_id    FROM   t_resource t1    JOIN t_monitor_product t3     ON t3.`name` = t1.product_name    WHERE   t1.tenant_id != ''    AND t3.abbreviation = ?"
	var sqlParam = []interface{}{productType}
	var tenantIds []string
	return util.Paginate(pageSize, pageCurrent, sql, sqlParam, &tenantIds)
}

func (mpd *AlarmInstanceDao) UpdateBatchInstanceName(models []*model.AlarmInstance) {
	sql := "UPDATE t_resource SET instance_name = CASE instance_id %s END WHERE instance_id IN ('%s')"
	var sql1 string
	var arr []string
	for _, v := range models {
		sql1 += " WHEN '" + v.InstanceID + "' THEN '" + v.InstanceName + "'"
		arr = append(arr, v.InstanceID)
	}
	sql2 := strings.Join(arr, "','")
	var i int
	global.DB.Raw(fmt.Sprintf(sql, sql1, sql2)).Find(&i)
}

func (mpd *AlarmInstanceDao) SelectInstanceList(tenantId string, productType string) []*model.AlarmInstance {
	sql := "SELECT DISTINCT instance_id  FROM t_resource t1, t_monitor_product t2  WHERE t1.tenant_id =?  AND t1.product_name = t2.NAME  and t2.abbreviation=? "
	var instance = &[]*model.AlarmInstance{}
	global.DB.Raw(sql, tenantId, productType).Find(instance)
	return *instance
}

func (mpd *AlarmInstanceDao) DeleteInstanceList(tenantId string, models []*model.AlarmInstance) {
	if len(models) == 0 {
		return
	}
	sql := "DELETE FROM t_resource WHERE tenant_id = ? and instance_id IN (?)"
	var ids []string
	for _, v := range models {
		ids = append(ids, v.InstanceID)
	}
	var i int
	db := global.DB
	db.Raw(sql, tenantId, ids).Find(&i)
	deleteRel := "delete FROM t_alarm_rule_resource_rel where tenant_id= ? and resource_id in (?)"
	db.Raw(deleteRel, tenantId, ids).Find(&i)
}
