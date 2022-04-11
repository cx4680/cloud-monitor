package main

import (
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/run_time"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq/consumer"
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
	loader := cp.NewMainLoader()
	loader.AddStage(func(*context.Context) error {
		cfg := config.Cfg.Iam
		middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
		return nil
	})
	loader.AddStage(func(*context.Context) error {
		return sys_db.InitData(config.Cfg.Db, "hawkeye", "file://./migrations")
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

	loader.AddStage(func(*context.Context) error {
		return k8s.InitK8s()
	})
	loader.AddStage(func(*context.Context) error {
		return sys_rocketmq.StartConsumersScribe(sys_rocketmq.Group(config.Cfg.Common.RegionName), []*sys_rocketmq.Consumer{{
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
			Handler: consumer.AlarmRecordAddHandler,
		}, {
			Topic:   sys_rocketmq.AlarmInfoTopic,
			Handler: consumer.AlarmInfoAddHandler,
		}, {
			Topic:   sys_rocketmq.ContactTopic,
			Handler: consumer.ContactHandler,
		}, {
			Topic:   sys_rocketmq.ContactGroupTopic,
			Handler: consumer.ContactGroupHandler,
		}, {
			Topic:   sys_rocketmq.RuleTopic,
			Handler: consumer.AlarmRuleHandler,
		}, {
			Topic:   sys_rocketmq.MonitorProductTopic,
			Handler: consumer.MonitorProductHandler,
		}, {
			Topic:   sys_rocketmq.MonitorItemTopic,
			Handler: consumer.MonitorItemHandler,
		},
		})
	})

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
