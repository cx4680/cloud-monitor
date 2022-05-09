package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"gorm.io/gorm"
)

var AlarmItem = new(AlarmItemDao)

type AlarmItemDao struct {
}

func (dao *AlarmItemDao) InsertBatch(db *gorm.DB, records []model.AlarmItem) {
	db.Create(&records)
}

func (dao *AlarmItemDao) GetItemListByRuleBizId(db *gorm.DB, ruleBizId string) []model.AlarmItem {
	var list []model.AlarmItem
	db.Model(&model.AlarmItem{}).Where("rule_biz_id=?", ruleBizId).Scan(&list)
	return list
}
