package data_sync

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/run_time"
	"errors"
	"time"
)

type Task interface {
	Run() error
}

type SyncTask interface {
	Task
	Period() time.Duration
	SetLoop(chan SyncTask)
	Loop() chan SyncTask
}

type CombinedSynchronizer struct {
	Loop  chan SyncTask
	Tasks []SyncTask
}

func NewCombinedSynchronizer(tasks []SyncTask) (Task, error) {
	if len(tasks) == 0 {
		return nil, errors.New("任务不能为空")
	}
	return &CombinedSynchronizer{
		Tasks: tasks,
		Loop:  make(chan SyncTask),
	}, nil
}

func (cs *CombinedSynchronizer) Run() error {
	go func() {
		for {
			select {
			case st := <-cs.Loop:
				go run_time.SafeRun(func() {
					err := st.Run()
					if err != nil {
						logger.Logger().Errorf("同步数据：出错原因：%v ", err)
					}
				})
			}
		}
	}()
	for i := range cs.Tasks {
		cs.Tasks[i].SetLoop(cs.Loop)
		cs.Loop <- cs.Tasks[i]
	}
	return nil
}
