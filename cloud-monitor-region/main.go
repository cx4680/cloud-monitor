package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/pipeline/sys_upgrade"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/middleware"
	"context"
	"fmt"
	"os"
)

func main() {

	loader := cp.NewMainLoader().
		AddStage(func(c *context.Context) error {
			cfg := config.Cfg.Iam
			middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
			return nil
		}).
		AddStage(func(c *context.Context) error {
			initializer := sys_db.DBInitializer{
				DB:      global.DB,
				Fetches: []sys_db.InitializerFetch{sys_db.CommonFetch, pipeline.ProjectFetch},
			}

			if err := initializer.Initnitialization(); err != nil {
				return err
			}
			return nil
		}).
		AddStage(func(c *context.Context) error {
			return translate.InitTrans("zh")
		}).
		AddStage(func(c *context.Context) error {
			bt := commonTask.NewBusinessTaskImpl()
			err := task.AddSyncJobs(bt)
			bt.Start()
			return err

		}).
		AddStage(func(c *context.Context) error {
			return k8s.InitK8s()
		}).
		AddStage(func(c *context.Context) error {
			return sys_rocketmq.StartConsumersScribe(sys_rocketmq.Group(config.Cfg.Common.RegionName), []*sys_rocketmq.Consumer{{
				Topic:   sys_rocketmq.AlertContactTopic,
				Handler: consumer.AlertContactHandler,
			}, {
				Topic:   sys_rocketmq.AlertContactGroupTopic,
				Handler: consumer.AlertContactGroupHandler,
			}, {
				Topic:   sys_rocketmq.RuleTopic,
				Handler: consumer.AlarmRuleHandler,
			}})
		}).
		AddStage(func(c *context.Context) error {
			sys_upgrade.PrometheusRuleUpgrade()
			return nil
		}).
		AddStage(func(c *context.Context) error {
			return web.Start(config.Cfg.Serve)
		})

	_, err := loader.Start()
	if err != nil {
		fmt.Printf("exit error:%v", err)
		os.Exit(1)
	}
}
