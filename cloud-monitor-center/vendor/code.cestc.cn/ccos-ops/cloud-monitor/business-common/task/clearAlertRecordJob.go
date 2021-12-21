package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
)

var clearIntervalDay = "180"

var name = "region-hawkeye-clearAlertRecordJob"

func Clear() {
	logger.Logger().Info("clearAlertRecordJob start")
	// TODO 锁
	dao.AlertRecordCommon.DeleteExpired(clearIntervalDay)
	logger.Logger().Info("clearAlertRecordJob end")
}
