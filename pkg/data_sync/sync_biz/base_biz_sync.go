package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
)

const (
	ContactSync     = "contact"
	AlarmSync       = "alarmRule"
	AlarmRecordSync = "alarmRecord"
)

type SyncDataTask interface {
	data_sync.SyncTask
	Work(string) (string, error)
}

type BaseSyncTask struct {
	Task    SyncDataTask
	BizCode string
}

func (bst *BaseSyncTask) Run() error {
	lastUpdateTime, err := bst.getUpdateTime()
	if err != nil {
		return err
	}
	currentTime, err := bst.Task.Work(lastUpdateTime)
	if err != nil {
		return err
	}
	return bst.updateTime(currentTime)
}

func (bst *BaseSyncTask) getUpdateTime() (string, error) {
	syncTime := dao.NewRegionSyncDao().GetUpdateTime(global.DB, bst.BizCode)
	return syncTime.UpdateTime, nil
}

func (bst *BaseSyncTask) updateTime(currentTime string) error {
	dao.NewRegionSyncDao().UpdateTime(global.DB, model.SyncTime{Name: bst.BizCode, UpdateTime: currentTime})
	return nil
}
