package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/task"
	"fmt"
	"os"
	"testing"
)

func Test111(t *testing.T) {
	os.Setenv("DB_PWD", "123456")
	config.InitConfig("D:\\dev-go\\cloud-monitor\\cloud-monitor-region\\config.local.yml")
	sys_db.InitDb(config.Cfg.Db)

	bt := commonTask.NewBusinessTaskImpl()
	go AddSyncJobs(bt)
	bt.Start()
	select {}
}

func TestRegion(t *testing.T) {
	os.Setenv("DB_PWD", "123456")
	config.InitConfig("D:\\dev-go\\cloud-monitor\\cloud-monitor-region\\config.local.yml")
	if err := sys_redis.InitClient(config.Cfg.Redis); err != nil {
		logger.Logger().Errorf("init redis error: %v\n", err)
	}
	d := GetRegionInfo("xx")
	fmt.Print(d)
}
