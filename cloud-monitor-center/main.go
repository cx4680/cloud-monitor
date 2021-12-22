package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
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

	loader := cp.NewMainLoader()
	loader.AddStage(func(c *context.Context) error {
		cfg := config.GetIamConfig()
		middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
		return nil
	})

	loader.AddStage(func(c *context.Context) error {
		initializer := sysDb.DBInitializer{
			DB:      global.DB,
			Fetches: []sysDb.InitializerFetch{new(sysDb.CommonInitializerFetch), new(pipeline.ProjectInitializerFetch)},
		}
		if err := initializer.Initnitialization(); err != nil {
			return err
		}
		return nil
	})

	loader.AddStage(func(c *context.Context) error {
		return translate.InitTrans("zh")
	})

	loader.AddStage(func(c *context.Context) error {
		var taskChan = make(chan error, 1)
		go func() {
			bt := task.NewBusinessTaskImpl()
			err := bt.Add(task.BusinessTaskDTO{
				Cron: "0 0 0/1 * * ?",
				Name: "clearAlertRecordJob",
				Task: task.Clear,
			})
			if err != nil {
				taskChan <- err
				return
			}

			bt.Start()
			close(taskChan)
		}()
		select {
		case err := <-taskChan:
			return err
		}
	})

	loader.AddStage(func(c *context.Context) error {
		return sysRocketMq.StartConsumersScribe("cloud-monitor-center", []*sysRocketMq.Consumer{{
			Topic:   sysRocketMq.InstanceTopic,
			Handler: consumer.InstanceHandler,
		}, {
			Topic:   sysRocketMq.SmsMarginReminderTopic,
			Handler: consumer.SmsMarginReminderConsumer,
		}, {
			Topic:   sysRocketMq.DeleteInstanceTopic,
			Handler: consumer.DeleteInstanceHandler,
		}, {
			Topic:   sysRocketMq.RecordTopic,
			Handler: consumer.AlertRecordAddHandler,
		}})
	})

	loader.AddStage(func(c *context.Context) error {
		return web.Start(config.GetServeConfig())
	})

	_, err := loader.Start()
	if err != nil {
		fmt.Printf("exit error: %v", err)
		os.Exit(1)
	}

}
