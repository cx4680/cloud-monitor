package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"gorm.io/gorm"
)

type AlertRecordCommonDao struct {
	db *gorm.DB
}

func NewAlertRecordCommonDao() *AlertRecordCommonDao {
	return &AlertRecordCommonDao{db: database.GetDb()}
}

func (mpd *AlertRecordCommonDao) DeleteExpired(day string) {
	mpd.db.Where("TO_DAYS(NOW()) - TO_DAYS(create_time) >= ?", day).Delete(models.AlertRecord{})
}
