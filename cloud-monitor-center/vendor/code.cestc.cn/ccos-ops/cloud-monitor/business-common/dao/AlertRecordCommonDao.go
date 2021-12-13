package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
)

type AlertRecordCommonDao struct {
}

var AlertRecordCommon = new(AlertRecordCommonDao)

func (mpd *AlertRecordCommonDao) DeleteExpired(day string) {
	global.DB.Where("TO_DAYS(NOW()) - TO_DAYS(create_time) >= ?", day).Delete(models.AlertRecord{})
}
