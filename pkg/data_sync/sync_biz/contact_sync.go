package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
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

func (cs *ContactSynchronizer) Work(lastUpdateTime string) error {

	return nil
}

func (cs *ContactSynchronizer) Period() time.Duration {
	return 10 * time.Second
}
