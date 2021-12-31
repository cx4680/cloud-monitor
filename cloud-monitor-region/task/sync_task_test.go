package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"testing"
)

func Test111(t *testing.T) {
	config.InitConfig("D:\\dev-go\\cloud-monitor\\cloud-monitor-region\\config.local.yml")
	sys_db.InitDb(config.GetDbConfig())
	bt := commonTask.NewBusinessTaskImpl()
	go AddSyncJobs(bt)
	bt.Start()
	select {}
}
