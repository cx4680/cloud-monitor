package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"gorm.io/gorm"
)

type InstanceDao struct {
	db *gorm.DB
}

func NewInstanceDao() *InstanceDao {
	return &InstanceDao{database.GetDb()}
}
func (dao *InstanceDao) GetInstanceNum(tenantId string) int {
	var result int
	dao.db.Raw(" SELECT count(DISTINCT ai.instance_id) num from t_alarm_instance ai       join t_alarm_rule ar on ar.id=ai.alarm_rule_id       where ar.tenant_id=?", tenantId).Scan(&result)
	return result
}
