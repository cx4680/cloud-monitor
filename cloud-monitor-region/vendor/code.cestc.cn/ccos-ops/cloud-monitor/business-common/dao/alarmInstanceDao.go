package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type AlarmInstanceDao struct {
	db *gorm.DB
}

func NewAlarmInstanceDao() *AlarmInstanceDao {
	return &AlarmInstanceDao{db: database.GetDb()}
}

func (mpd *AlarmInstanceDao) CountRegionInstanceNum(tenantId string) {

}

func (mpd *AlarmInstanceDao) SelectTenantIdList(productType string, pageCurrent int, pageSize int) []string {
	sql := "SELECT DISTINCT t1.tenant_id FROM t_alarm_instance t1 LEFT JOIN t_alarm_rule t2 ON t1.alarm_rule_id = t2.id WHERE t1.tenant_id != '' AND t2.product_type = '%s' LIMIT %s,%s"
	var tenantIds []string
	mpd.db.Raw(fmt.Sprintf(sql, productType, strconv.Itoa((pageCurrent-1)*pageSize), strconv.Itoa(pageSize))).Find(tenantIds)
	return tenantIds
}

func (mpd *AlarmInstanceDao) UpdateBatchInstanceName(models []models.AlarmInstance) {
	sql := "UPDATE t_alarm_instance SET instance_name = CASE instance_id %s END WHERE instance_id IN ('%s')"
	var sql1 string
	var arr []string
	for _, v := range models {
		sql1 += " WHEN '" + v.InstanceID + "' THEN '" + v.InstanceName + "'"
		arr = append(arr, v.InstanceID)
	}
	sql2 := strings.Join(arr, "','")
	var i int
	mpd.db.Raw(fmt.Sprintf(sql, sql1, sql2)).Find(i)
}

func (mpd *AlarmInstanceDao) SelectInstanceList(tenantId string, productType string) []models.AlarmInstance {
	sql := "SELECT t1.* FROM t_alarm_instance t1 LEFT JOIN t_alarm_rule t2 ON t1.alarm_rule_id = t2.id where t1.tenant_id = '%s' and t2.product_type = '%s'"
	var model = &[]models.AlarmInstance{}
	mpd.db.Raw(fmt.Sprintf(sql, tenantId, productType)).Find(model)
	return *model
}

func (mpd *AlarmInstanceDao) DeleteInstanceList(tenantId string, models []models.AlarmInstance) {
	sql := "DELETE FROM t_alarm_instance WHERE tenant_id = %s and instance_id IN ('%s')"
	var arr []string
	for _, v := range models {
		arr = append(arr, v.InstanceID)
	}
	sql1 := strings.Join(arr, "','")
	var i int
	mpd.db.Raw(fmt.Sprintf(sql, tenantId, sql1)).Find(i)
}
