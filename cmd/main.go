package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/run_time"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/data_sync/sync_init"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/pipeline/sys_upgrade"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/web"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/middleware"
	"context"
	"fmt"
	"os"
)

// @title Swagger Example API
// @version 0.0.1
// @description  This is a sample server Petstore server.
// @BasePath /
func main() {
	loader := pipeline.NewMainLoader()
	loader.AddStage(func(*context.Context) error {
		cfg := config.Cfg.Iam
		middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
		return nil
	})

	loader.AddStage(func(*context.Context) error {
		if err := sys_db.InitDb(config.Cfg.Db); err != nil {
			logger.Logger().Errorf("init database error: %v\n", err)
			//return err
		}
		return nil
	})

	loader.AddStage(func(*context.Context) error {
		if err := sys_redis.InitClient(config.Cfg.Redis); err != nil {
			logger.Logger().Errorf("init redis error: %v\n", err)
			//return err
		}
		return nil
	})

	//loader.AddStage(func(*context.Context) error {
	//	return sys_db.InitData(config.Cfg.Db, "hawkeye", "file://./migrations")
	//})

	loader.AddStage(func(c *context.Context) error {
		return sync_init.InitSync(config.Cfg.Common.RegionRole)
	})

	loader.AddStage(func(*context.Context) error {
		return translate.InitTrans("zh")
	})

	loader.AddStage(func(*context.Context) error {
		bt := task.NewBusinessTaskImpl()
		err := bt.Add(task.BusinessTaskDTO{
			Cron: "0 0 0/1 * * ?",
			Name: "clearAlarmRecordJob",
			Task: task.Clear,
		})
		err = task.AddSyncJobs(bt)
		bt.Start()
		return err
	})

	//loader.AddStage(func(*context.Context) error {
	//	return k8s.InitK8s()
	//})

	//loader.AddStage(func(c *context.Context) error {
	//	levelList := dao.ConfigItem.GetConfigItemList(dao.AlarmLevel)
	//	ls := make([]uint8, len(levelList))
	//	for i, l := range levelList {
	//		v, err := strconv.Atoi(l.Code)
	//		if err != nil {
	//			return err
	//		}
	//		ls[i] = uint8(v)
	//	}
	//	_ = k8s.DeleteAlertManagerConfig(k8s.LevelInhibitName)
	//
	//	return k8s.ApplyInhibitRules(ls)
	//})

	loader.AddStage(func(*context.Context) error {
		sys_upgrade.PrometheusRuleUpgrade()
		return nil
	})

	loader.AddStage(func(c *context.Context) error {
		go run_time.SafeRun(func() {
			service.StartHandleAlarmEvent(*c)
		})
		return nil
	})

	loader.AddStage(func(*context.Context) error {
		return web.Start(config.Cfg.Serve)
	})

	_, err := loader.Start()
	if err != nil {
		fmt.Printf("exit error: %v", err)
		os.Exit(1)
	}

}
