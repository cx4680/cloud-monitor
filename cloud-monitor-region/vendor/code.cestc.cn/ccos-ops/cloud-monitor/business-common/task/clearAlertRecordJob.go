package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"log"
)

var clearIntervalDay = "180"

var name = "region-hawkeye-clearAlertRecordJob"

func Clear() {
	log.Println("clearAlertRecordJob start")
	// TODO 锁
	dao.AlertRecordCommon.DeleteExpired(clearIntervalDay)
	log.Println("clearAlertRecordJob end")
}
