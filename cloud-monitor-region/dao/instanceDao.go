package dao

import "github.com/jinzhu/gorm"

type InstanceDao struct {
	db *gorm.DB
}

func NewInstanceDao(db *gorm.DB) *InstanceDao {
	return &InstanceDao{db}
}
func (dao *InstanceDao) GetInstanceNum(tenantId string) int {
	var result int
	dao.db.Raw(" SELECT count(DISTINCT ai.instance_id) num from t_alarm_instance ai       join t_alarm_rule ar on ar.id=ai.alarm_rule_id       where ar.tenant_id=?", tenantId).Scan(&result)
	return result
}
