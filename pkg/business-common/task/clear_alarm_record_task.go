package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
)

var clearIntervalDay = "180"

func Clear() {
	logger.Logger().Info("clearAlertRecordJob start")
	// TODO 锁
	dao.AlarmRecordCommon.DeleteExpired(clearIntervalDay)
	logger.Logger().Info("clearAlertRecordJob end")
}
