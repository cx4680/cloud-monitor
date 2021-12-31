package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"log"
	"testing"
)

func TestAlarmHandlerDao_GetHandlerListByRuleId(t *testing.T) {
	config.InitConfig("config.local.yml")
	sys_db.InitDb(config.GetDbConfig())
	ruleId := "1"
	list := AlarmHandler.GetHandlerListByRuleId(global.DB, ruleId)
	for _, handler := range list {
		log.Println(jsonutil.ToString(handler))
	}
}
