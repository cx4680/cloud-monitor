package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_db"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestAlertRecordDao_FindContactInfoByGroupIds(t *testing.T) {
	config.InitConfig("config.local.yml")
	sys_db.InitDb(config.Cfg.Db)
	list := AlarmRecord.FindContactInfoByGroupIds([]string{"1"})
	for _, info := range list {
		log.Println(info)
	}
}

func TestFloat(t *testing.T) {
	var f float64 = 100000000.2
	print(fmt.Sprintf("%.f", f))
}

func TestFindGroupIdsByRecordId(t *testing.T) {
	config.InitConfig("C:\\work\\go-space\\cloud-monitor\\cloud-monitor-center\\config.local.yml")
	os.Setenv("DB_PWD", "123456")
	sys_db.InitDb(config.Cfg.Db)
}
