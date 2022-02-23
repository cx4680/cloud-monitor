package main

import (
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/pipeline/sys_upgrade"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/run_time"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/middleware"
	"context"
	"os"
)

func main() {

	loader := cp.NewMainLoader().
		AddStage(func(*context.Context) error {
			cfg := config.Cfg.Iam
			middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
			return nil
		}).
		AddStage(func(*context.Context) error {
			return sys_db.InitData(config.Cfg.Db, "hawkeye_region", "file://./migrations")
		}).
		AddStage(func(*context.Context) error {
			return translate.InitTrans("zh")
		}).
		AddStage(func(*context.Context) error {
			bt := commonTask.NewBusinessTaskImpl()
			err := task.AddSyncJobs(bt)
			bt.Start()
			return err

		}).
		AddStage(func(*context.Context) error {
			return k8s.InitK8s()
		}).
		AddStage(func(*context.Context) error {
			return sys_rocketmq.StartConsumersScribe(sys_rocketmq.Group(config.Cfg.Common.RegionName), []*sys_rocketmq.Consumer{{
				Topic:   sys_rocketmq.ContactTopic,
				Handler: consumer.ContactHandler,
			}, {
				Topic:   sys_rocketmq.ContactGroupTopic,
				Handler: consumer.ContactGroupHandler,
			}, {
				Topic:   sys_rocketmq.RuleTopic,
				Handler: consumer.AlarmRuleHandler,
			}})
		}).
		AddStage(func(*context.Context) error {
			sys_upgrade.PrometheusRuleUpgrade()
			return nil
		}).
		AddStage(func(c *context.Context) error {
			go run_time.SafeRun(func() {
				service.StartHandleAlarmEvent(*c)
			})
			return nil
		}).
		AddStage(func(*context.Context) error {
			return web.Start(config.Cfg.Serve)
		})

	_, err := loader.Start()
	if err != nil {
		logger.Logger().Error("exit error, ", err)
		os.Exit(1)
	}
}
