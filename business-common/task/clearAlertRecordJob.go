package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"github.com/robfig/cron"
	"log"
)

var clearIntervalDay = "180"

var name = "region-hawkeye-clearAlertRecordJob"

func CronClear() {
	c := cron.New()
	err := c.AddFunc("0 0 0/1 * * ?", func() {
		clear()
	})
	if err != nil {
		log.Println("clearAlertRecordJob error", err)
	}
	c.Start()
	defer c.Stop()
	select {}
}

func clear() {
	log.Println("clearAlertRecordJob start")
	// TODO 锁
	dao.NewAlertRecordCommonDao(database.GetDb()).DeleteExpired(clearIntervalDay)
	log.Println("clearAlertRecordJob end")
}
