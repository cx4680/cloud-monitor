package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/pipeline"
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
			cfg := config.GetIamConfig()
			middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
			return nil
		}).
		AddStage(func(c *context.Context) error {
			initializer := sysDb.DBInitializer{
				DB:      global.DB,
				Fetches: []sysDb.InitializerFetch{new(sysDb.CommonInitializerFetch), new(pipeline.ProjectInitializerFetch)},
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
			var taskChan = make(chan error, 1)
			go func() {
				bt := commonTask.NewBusinessTaskImpl()
				err := task.AddSyncJobs(bt)
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
		}).
		AddStage(func(c *context.Context) error {
			return k8s.InitK8s()
		}).
		AddStage(func(c *context.Context) error {
			return sysRocketMq.StartConsumersScribe(sysRocketMq.Group(config.GetCommonConfig().RegionName), []*sysRocketMq.Consumer{{
				Topic:   sysRocketMq.AlertContactTopic,
				Handler: consumer.AlertContactHandler,
			}, {
				Topic:   sysRocketMq.AlertContactGroupTopic,
				Handler: consumer.AlertContactGroupHandler,
			}, {
				Topic:   sysRocketMq.RuleTopic,
				Handler: consumer.AlarmRuleHandler,
			}})
		}).
		AddStage(func(c *context.Context) error {
			return web.Start(config.GetServeConfig())
		})

	_, err := loader.Start()
	if err != nil {
		fmt.Printf("exit error:%v", err)
		os.Exit(1)
	}
}
