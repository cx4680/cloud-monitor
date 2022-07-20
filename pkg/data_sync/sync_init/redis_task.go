package sync_init

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/run_time"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type MyExclusiveTask struct {
	Id       string
	Name     string
	MasterId string
	Ready    chan int
	Do       func(isMaster func() bool) error
}

func NewMyExclusiveTask(name, id string, work func(isMaster func() bool) error) *MyExclusiveTask {
	t := &MyExclusiveTask{
		Id:   id,
		Name: name,
		Do:   work,
	}
	t.Ready = make(chan int, 1)

	return t
}

func (met *MyExclusiveTask) isMaster() bool {
	return met.Id == met.MasterId
}

func (met *MyExclusiveTask) Run() error {
	defer func() {
		close(met.Ready)
	}()

	if err := met.checkMasterSurviveTask(); err != nil {
		return err
	}

	if err := met.start(); err != nil {
		return err
	}
	return nil
}

func (met *MyExclusiveTask) getTaskStartedFlagKey() string {
	return fmt.Sprintf(constant.SyncTaskKey, met.Name)
}

func (met *MyExclusiveTask) getMasterId() (string, error) {
	v, err := sys_redis.Get(met.getTaskStartedFlagKey())
	if err != nil && err != redis.Nil {
		return "", err
	}
	if len(v) == 0 {
		return "", nil
	}
	return v, nil
}

func (met *MyExclusiveTask) start() error {
	ctx := context.Background()
	lockKey := fmt.Sprintf(constant.SyncTaskStartLockKey, met.Name)

	if err := sys_redis.Lock(ctx, lockKey, 10*time.Second, true); err != nil {
		logger.Logger().Info(met.Id + " get lock fail")
		return err
	}

	defer func() {
		if err := sys_redis.Unlock(ctx, lockKey); err != nil {
			logger.Logger().Error(err.Error())
		}
		logger.Logger().Info(met.Id + " release lock.")
	}()

	logger.Logger().Info(met.Id + " get lock success!")

	masterId, err := met.getMasterId()
	if err != nil {
		return err
	}

	if len(masterId) > 0 {
		met.MasterId = masterId
		logger.Logger().Info(met.Id + " get result: " + met.Name + " already started, masterId=" + met.MasterId)
		return nil
	}
	logger.Logger().Info(met.Id + " starting...")
	err = met.masterHeartBeatTask(met.isMaster)
	if err != nil {
		return err
	}
	met.MasterId = met.Id
	return met.Do(met.isMaster)
}

func (met *MyExclusiveTask) checkMasterSurviveTask() error {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		select {
		case <-met.Ready:
			for {
				select {
				case <-ticker.C:
					run_time.SafeRun(func() {
						v, err := sys_redis.Get(met.getTaskStartedFlagKey())
						logger.Logger().Info(met.Id + " check master survive， masterId=" + v)
						if err != nil && err != redis.Nil {
							return
						}
						if len(v) == 0 {
							logger.Logger().Info(met.Id + " found master changed...")
							if e := met.start(); e != nil {
								logger.Logger().Error(e.Error())
							}
						}
					})
				}
			}
		}
	}()
	return nil
}

func (met *MyExclusiveTask) masterHeartBeatTask(isMaster func() bool) error {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer func() {
			ticker.Stop()
		}()

		op := func() {
			logger.Logger().Info(met.Id + " heartbeat ")
			err := sys_redis.SetByTimeOut(met.getTaskStartedFlagKey(), met.Id, 30*time.Second)
			if err != nil {
				logger.Logger().Error(err.Error())
			}
		}
		select {
		case <-met.Ready:
			op()
		}
		for {
			select {
			case <-ticker.C:
				if isMaster() {
					run_time.SafeRun(op)
				}
			}
		}
	}()
	return nil
}
