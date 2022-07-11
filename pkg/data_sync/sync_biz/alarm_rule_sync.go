package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
	"log"
	"time"
)

type AlarmRuleSynchronizer struct {
	BaseSyncTask
}

func NewAlarmRuleSynchronizer() data_sync.SyncTask {
	s := &AlarmRuleSynchronizer{
		BaseSyncTask: BaseSyncTask{
			BizCode: AlarmSync,
		},
	}
	s.Task = s
	return s
}

func (ars *AlarmRuleSynchronizer) Work(lastUpdateTime string) error {
	log.Println("告警规则同步...")
	time.Sleep(10 * time.Second)
	return nil
}

func (ars *AlarmRuleSynchronizer) Period() time.Duration {
	return 30 * time.Second
}
