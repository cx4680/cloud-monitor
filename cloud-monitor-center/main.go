package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
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
		cfg := config.Cfg.Iam
		middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
		return nil
	})

	loader.AddStage(func(c *context.Context) error {
		initializer := sys_db.DBInitializer{
			DB:      global.DB,
			Fetches: []sys_db.InitializerFetch{new(sys_db.CommonInitializerFetch), new(pipeline.ProjectInitializerFetch)},
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
		bt := task.NewBusinessTaskImpl()
		err := bt.Add(task.BusinessTaskDTO{
			Cron: "0 0 0/1 * * ?",
			Name: "clearAlertRecordJob",
			Task: task.Clear,
		})
		bt.Start()
		return err
	})

	loader.AddStage(func(c *context.Context) error {
		return sys_rocketmq.StartConsumersScribe("cloud-monitor-center", []*sys_rocketmq.Consumer{{
			Topic:   sys_rocketmq.InstanceTopic,
			Handler: consumer.InstanceHandler,
		}, {
			Topic:   sys_rocketmq.SmsMarginReminderTopic,
			Handler: consumer.SmsMarginReminderConsumer,
		}, {
			Topic:   sys_rocketmq.DeleteInstanceTopic,
			Handler: consumer.DeleteInstanceHandler,
		}, {
			Topic:   sys_rocketmq.RecordTopic,
			Handler: consumer.AlertRecordAddHandler,
		}})
	})

	loader.AddStage(func(c *context.Context) error {
		return web.Start(config.Cfg.Serve)
	})

	_, err := loader.Start()
	if err != nil {
		fmt.Printf("exit error: %v", err)
		os.Exit(1)
	}

}
