package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
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

func (ars *AlarmRuleSynchronizer) Work(lastUpdateTime string) (string, error) {
	currentTime, err := service.NewRegionSyncService().AlarmRuleSync(lastUpdateTime)
	time.Sleep(10 * time.Second)
	return currentTime, err
}

func (ars *AlarmRuleSynchronizer) Period() time.Duration {
	return 30 * time.Second
}
