package sync_init

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/run_time"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync/sync_biz"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
	"time"
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

	//同步器心跳检测
	go run_time.SafeRun(func() {
		startHeartBeatCheckTask()
	})
	return startSyncTask()

}

func startHeartBeatCheckTask() {
	logger.Logger().Info("start sync heartbeat check ")
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			val, err := sys_redis.Get(constant.SyncFlagKey)
			if err != nil && err != redis.Nil {
				logger.Logger().Infof("get started flag error, %v", err)
			}
			if len(val) == 0 {
				logger.Logger().Info("sync task is null, start running sync task")
				err := startSyncTask()
				if err != nil {
					logger.Logger().Info("start sync task error, %v", err)
				}
			}
		}
	}
}

func startSyncTask() error {
	ctx := context.Background()
	lockKey := constant.SyncStartKey
	var started = make(chan int)

	defer func() {
		//休眠1s，确保已成功发送心跳
		<-started
		err := sys_redis.Unlock(ctx, lockKey)
		logger.Logger().Info("release sync lock")
		if err != nil {
			logger.Logger().Errorf("unlock errorL%+v, lock:%s", err, lockKey)
		}
	}()

	err := sys_redis.Lock(ctx, lockKey, sys_redis.DefaultLease, true)
	if err != nil {
		log.Printf("获取 rule lock error  %+v", err)
		return err
	}
	logger.Logger().Info("get SyncStartKey lock")
	val, err := sys_redis.Get(constant.SyncFlagKey)
	if err != nil && err != redis.Nil {
		logger.Logger().Errorf("get sync flag key error, %v", err)
		return nil
	}
	if len(val) > 0 {
		logger.Logger().Info("already start sync task, not need running ")
		started <- 1
		return nil
	}
	logger.Logger().Info("sync task not running, start running...")
	synchronizer, err := data_sync.NewCombinedSynchronizer([]data_sync.SyncTask{sync_biz.NewContactSynchronizer(), sync_biz.NewAlarmRuleSynchronizer(), sync_biz.NewAlarmRecordSynchronizer()})
	if err != nil {
		return err
	}
	err = synchronizer.Run()
	if err != nil {
		return err
	}
	heartBeatTask()
	started <- 1
	logger.Logger().Info("sync task is running success!")
	return nil
}

func heartBeatTask() {
	heartBeatFlag := func() {
		err := sys_redis.SetByTimeOut(constant.SyncFlagKey, "started", 30*time.Second)
		if err != nil {
			logger.Logger().Info("set started flag error, %v", err)
		}
	}
	heartBeatFlag()
	//同步器心跳
	go run_time.SafeRun(func() {
		logger.Logger().Info("start sync heartbeat... ")
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				heartBeatFlag()
			}
		}
	})
}
