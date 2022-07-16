package sync_init

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync/sync_biz"
	"github.com/google/uuid"
	"strings"
)

func InitSync(regionRole string) error {
	if strings.EqualFold(regionRole, "integration") {
		logger.Logger().Info("current run mode: single region")
		return nil
	}
	if strings.EqualFold(regionRole, "manager") {
		logger.Logger().Info("current run mode: multiple regions - manager")
		return nil
	}
	logger.Logger().Info("current run mode: multiple regions - region")

	err := NewMyExclusiveTask("syncDataTask", uuid.New().String(), func(isMaster func() bool) error {
		synchronizer, err := data_sync.NewCombinedSynchronizer([]data_sync.SyncTask{sync_biz.NewContactSynchronizer(),
			sync_biz.NewAlarmRuleSynchronizer(), sync_biz.NewAlarmRecordSynchronizer()})
		if err != nil {
			return err
		}
		return synchronizer.Run(isMaster)
	}).Run()

	return err

}
