package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"gorm.io/gorm"
)

type AlertRecordCommonDao struct {
	db *gorm.DB
}

func NewAlertRecordCommonDao(db *gorm.DB) *AlertRecordCommonDao {
	return &AlertRecordCommonDao{db: db}
}

func (mpd *AlertRecordCommonDao) DeleteExpired(day string) {
	mpd.db.Where("TO_DAYS(NOW()) - TO_DAYS(create_time) >= ?", day).Delete(models.AlertRecord{})
}
