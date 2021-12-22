package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"log"
	"testing"
)

func TestAlertRecordDao_FindContactInfoByGroupIds(t *testing.T) {
	config.InitConfig("config.local.yml")
	sysDb.InitDb(config.GetDbConfig())
	list := AlertRecord.FindContactInfoByGroupIds([]string{"1"})
	for _, info := range list {
		log.Println(info)
	}
}
