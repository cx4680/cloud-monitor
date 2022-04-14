package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
)

type AlarmRecordCommonDao struct {
}

var AlarmRecordCommon = new(AlarmRecordCommonDao)

func (mpd *AlarmRecordCommonDao) DeleteExpired(day string) {
	global.DB.Where("TO_DAYS(NOW()) - TO_DAYS(create_time) >= ?", day).Delete(model.AlarmRecord{})
}
