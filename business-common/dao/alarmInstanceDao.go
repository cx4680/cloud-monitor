package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/pageUtils"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"fmt"
	"strings"
	"unsafe"
)

type AlarmInstanceDao struct {
}

var AlarmInstance = new(AlarmInstanceDao)

func (mpd *AlarmInstanceDao) CountRegionInstanceNum(tenantId string) {

}

func (mpd *AlarmInstanceDao) SelectTenantIdList(productType string, pageCurrent int, pageSize int) *vo.PageVO {
	sql := "SELECT DISTINCT   t1.tenant_id FROM   t_alarm_instance t1 LEFT JOIN t_alarm_rule t2 ON t1.alarm_rule_id = t2.id LEFT JOIN monitor_product t3 ON t3.`name` = t2.product_type WHERE   t1.tenant_id != '' AND t3.id  = ?"
	var sqlParam = []interface{}{productType}
	var tenantIds []string
	return pageUtils.Paginate(pageSize, pageCurrent, sql, sqlParam, unsafe.Pointer(&tenantIds))
}

func (mpd *AlarmInstanceDao) UpdateBatchInstanceName(models []*models.AlarmInstance) {
	sql := "UPDATE t_alarm_instance SET instance_name = CASE instance_id %s END WHERE instance_id IN ('%s')"
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

func (mpd *AlarmInstanceDao) SelectInstanceList(tenantId string, productType string) []*models.AlarmInstance {
	sql := "SELECT   t1.* FROM     t_alarm_instance t1    LEFT JOIN t_alarm_rule t2 ON t1.alarm_rule_id = t2.id     LEFT JOIN monitor_product t3 on t3.`name`=t2.product_type     WHERE    t1.tenant_id = ?    AND t3.id =?"
	var model = &[]*models.AlarmInstance{}
	global.DB.Raw(sql, tenantId, productType).Find(model)
	return *model
}

func (mpd *AlarmInstanceDao) DeleteInstanceList(tenantId string, models []*models.AlarmInstance) {
	sql := "DELETE FROM t_alarm_instance WHERE tenant_id = %s and instance_id IN ('%s')"
	var arr []string
	for _, v := range models {
		arr = append(arr, v.InstanceID)
	}
	sql1 := strings.Join(arr, "','")
	var i int
	global.DB.Raw(fmt.Sprintf(sql, tenantId, sql1)).Find(i)
}
