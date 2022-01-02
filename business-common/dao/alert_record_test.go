package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"fmt"
	"log"
	"testing"
)

func TestAlertRecordDao_FindContactInfoByGroupIds(t *testing.T) {
	config.InitConfig("config.local.yml")
	sys_db.InitDb(config.Cfg.Db)
	list := AlertRecord.FindContactInfoByGroupIds([]string{"1"})
	for _, info := range list {
		log.Println(info)
	}
}

func TestFloat(t *testing.T) {
	var f float64 = 100000000.2
	print(fmt.Sprintf("%.f", f))
}
