package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_db"
	"log"
	"testing"
)

func TestAlarmHandlerDao_GetHandlerListByRuleId(t *testing.T) {
	config.InitConfig("config.local.yml")
	sys_db.InitDb(config.Cfg.Db)
	ruleId := "1"
	list := AlarmHandler.GetHandlerListByRuleId(global.DB, ruleId)
	for _, handler := range list {
		log.Println(jsonutil.ToString(handler))
	}
}
