package sync

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"strings"
	"time"
)

//const defaultEmptyNameSrv = "${rocketMQEndpoint}"

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
	contact := AlarmContact{}
	execute(contact)
	//rule := AlarmRule{}
	//execute(rule)
	//history := AlarmHistory{}
	//execute(history)
	return nil
}

type SyncCornTasker interface {
	SyncExecute()
}

type AlarmContact struct{}

type AlarmRule struct{}

type AlarmHistory struct{}

func (c AlarmContact) SyncExecute() {
	err := service.NewRegionSyncService().ContactSync()
	if err != nil {
		logger.Logger().Errorf("同步失败：%v", err)
	}
}

func (r AlarmRule) SyncExecute() {
	err := service.NewRegionSyncService().AlarmRuleSync()
	if err != nil {
		logger.Logger().Errorf("同步失败：%v", err)
	}
}

func (h AlarmHistory) SyncExecute() {
	err := service.NewRegionSyncService().AlarmRecordSync()
	if err != nil {
		logger.Logger().Errorf("同步失败：%v", err)
	}
}

func execute(arg SyncCornTasker) {
	loop := make(chan int, 1)
	sync := func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger().Info("同步数据：失败原因：v%\n", err)
			}
		}()
		arg.SyncExecute()
		time.Sleep(30 * time.Second)
		loop <- 1
	}
	loop <- 1
	go func() {
		for {
			select {
			case <-loop:
				sync()
			}
		}
	}()
}
