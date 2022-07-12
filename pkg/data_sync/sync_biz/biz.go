package sync_biz

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
	"time"
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
	loop    chan data_sync.SyncTask
}

func (bst *BaseSyncTask) Period() time.Duration {
	return bst.Task.Period()
}
func (bst *BaseSyncTask) SetLoop(loop chan data_sync.SyncTask) {
	bst.loop = loop
}
func (bst *BaseSyncTask) Loop() chan data_sync.SyncTask {
	return bst.loop
}

func (bst *BaseSyncTask) Work(time string) (string, error) {
	currentTime, err := bst.Task.Work(time)
	return currentTime, err
}

func (bst *BaseSyncTask) Run() error {
	defer func() {
		time.Sleep(bst.Period())
		bst.Loop() <- bst.Task
	}()
	lastUpdateTime, err := bst.getUpdateTime()
	if err != nil {
		return err
	}
	currentTime, err := bst.Task.Work(lastUpdateTime)
	if err != nil {
		return err
	}
	err = bst.updateTime(currentTime)
	if err != nil {
		return err
	}
	return nil
}

func (bst *BaseSyncTask) getUpdateTime() (string, error) {
	syncTime := dao.NewRegionSyncDao().GetUpdateTime(global.DB, bst.BizCode)
	return syncTime.UpdateTime, nil
}

func (bst *BaseSyncTask) updateTime(currentTime string) error {
	dao.NewRegionSyncDao().UpdateTime(global.DB, model.SyncTime{Name: bst.BizCode, UpdateTime: currentTime})
	return nil
}
