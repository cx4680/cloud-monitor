package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"gorm.io/gorm"
)

type AlarmHandlerDao struct {
}

var AlarmHandler = new(AlarmHandlerDao)

func (dao *AlarmHandlerDao) GetHandlerListByRuleId(db *gorm.DB, ruleId string) []models.AlarmHandler {
	var list []models.AlarmHandler
	db.Where(models.AlarmHandler{
		AlarmRuleId: ruleId,
	}).Find(&list)
	return list
}
