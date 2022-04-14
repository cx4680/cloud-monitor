package dto

import (
	model2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
)

type AlarmSyncData struct {
	RecordList []model2.AlarmRecord
	InfoList   []model2.AlarmInfo
}
