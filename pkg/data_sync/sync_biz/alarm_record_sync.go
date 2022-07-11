package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
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

func (ars *AlarmRecordSynchronizer) Work(lastUpdateTime string) error {

	return nil

}

func (ars *AlarmRecordSynchronizer) Period() time.Duration {
	return 30 * time.Second
}
