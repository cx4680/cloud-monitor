package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"gorm.io/gorm"
)

type AlarmHandlerDao struct {
}

var AlarmHandler = new(AlarmHandlerDao)

func (dao *AlarmHandlerDao) GetHandlerListByRuleId(db *gorm.DB, ruleId string) []model.AlarmHandler {
	var list []model.AlarmHandler
	db.Where(model.AlarmHandler{
		AlarmRuleId: ruleId,
	}).Find(&list)
	return list
}
