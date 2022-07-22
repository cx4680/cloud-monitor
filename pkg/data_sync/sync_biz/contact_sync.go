package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"time"
)

type ContactSynchronizer struct {
	BaseSyncTask
}

func NewContactSynchronizer() data_sync.SyncTask {
	s := &ContactSynchronizer{
		BaseSyncTask: BaseSyncTask{
			BizCode: ContactSync,
		},
	}
	s.Task = s
	return s
}

func (cs *ContactSynchronizer) Work(lastUpdateTime string) (string, error) {
	currentTime, err := service.NewRegionSyncService().ContactSync(lastUpdateTime)
	return currentTime, err
}

func (cs *ContactSynchronizer) Period() time.Duration {
	return 30 * time.Second
}
