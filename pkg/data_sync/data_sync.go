package data_sync

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"errors"
	"time"
)

type Task interface {
	Run() error
}
type ClusterTask interface {
	Run(isMaster func() bool) error
}

type SyncTask interface {
	Task
	Period() time.Duration
}

type CombinedSynchronizer struct {
	Loop  chan SyncTask
	Tasks []SyncTask
}

func NewCombinedSynchronizer(tasks []SyncTask) (ClusterTask, error) {
	if len(tasks) == 0 {
		return nil, errors.New("任务不能为空")
	}
	return &CombinedSynchronizer{
		Tasks: tasks,
		Loop:  make(chan SyncTask, 10),
	}, nil
}

func (cs *CombinedSynchronizer) Run(isMaster func() bool) error {
	go func() {
		for {
			select {
			case st := <-cs.Loop:
				if !isMaster() {
					goto EXIT
				}
				go func(t SyncTask) {
					defer func() {
						if e := recover(); e != nil {
							logger.Logger().Errorf("sync run time error, %v", e)
						}
						time.Sleep(st.Period())
						cs.Loop <- st
					}()
					err := st.Run()
					if err != nil {
						logger.Logger().Errorf("同步数据：出错原因：%v ", err)
					}
				}(st)
			}
		}
	EXIT:
		logger.Logger().Info("master killed...")
	}()
	for i := range cs.Tasks {
		cs.Loop <- cs.Tasks[i]
	}

	return nil
}
