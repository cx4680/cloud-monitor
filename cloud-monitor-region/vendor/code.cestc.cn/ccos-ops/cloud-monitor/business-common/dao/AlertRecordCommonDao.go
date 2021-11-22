package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
)

type AlertRecordCommonDao struct {
}

var AlertRecordCommon = new(AlertRecordCommonDao)

func (mpd *AlertRecordCommonDao) DeleteExpired(day string) {
	database.GetDb().Where("TO_DAYS(NOW()) - TO_DAYS(create_time) >= ?", day).Delete(models.AlertRecord{})
}
