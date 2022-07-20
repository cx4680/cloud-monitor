package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"time"
)

type AlarmRecordSynchronizer struct {
	BaseSyncTask
}

func NewAlarmRecordSynchronizer() data_sync.SyncTask {
	s := &AlarmRecordSynchronizer{
		BaseSyncTask: BaseSyncTask{
			BizCode: AlarmRecordSync,
		},
	}
	s.Task = s
	return s
}

func (ars *AlarmRecordSynchronizer) Work(lastUpdateTime string) (string, error) {
	currentTime, err := service.NewRegionSyncService().AlarmRecordSync(lastUpdateTime)
	return currentTime, err
}

func (ars *AlarmRecordSynchronizer) Period() time.Duration {
	return 10 * time.Second
}
