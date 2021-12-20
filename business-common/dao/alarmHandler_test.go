package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"log"
	"testing"
)

func TestAlarmHandlerDao_GetHandlerListByRuleId(t *testing.T) {
	config.InitConfig("config.local.yml")
	sysDb.InitDb(config.GetDbConfig())
	ruleId := "1"
	list := AlarmHandler.GetHandlerListByRuleId(global.DB, ruleId)
	for _, handler := range list {
		log.Println(tools.ToString(handler))
	}
}
