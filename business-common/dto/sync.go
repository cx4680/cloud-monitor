package dto

import "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"

type AlarmSyncData struct {
	RecordList []model.AlarmRecord
	InfoList   []model.AlarmInfo
}
