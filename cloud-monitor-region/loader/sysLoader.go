package loader

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"log"
)

type TransactionLoader struct {
}

func (t *TransactionLoader) Load() error {
	return translate.InitTrans("zh")
}

type TaskLoader struct{}

func (t *TaskLoader) Load() error {
	bt := commonTask.NewBusinessTaskImpl()
	if err := bt.Add(commonTask.BusinessTaskDTO{
		Cron: "0 0 0/1 * * ?",
		Name: "instanceJob",
		Task: task.NewEcsJob().SyncJob,
	}); err != nil {
		return err
	}

	if err := bt.Add(commonTask.BusinessTaskDTO{
		Cron: "0 0 0/1 * * ?",
		Task: task.NewSlbJob().SyncJob,
	}); err != nil {
		return err
	}
	if err := bt.Add(commonTask.BusinessTaskDTO{
		Cron: "0 0 0/1 * * ?",
		Task: task.NewEipJob().SyncJob,
	}); err != nil {
		return err
	}

	bt.Start()
	return nil
}

type RocketMQConsumerLoader struct{}

func (r *RocketMQConsumerLoader) Load() error {
	if err := sysRocketMq.StartConsumersScribe(sysRocketMq.Group(config.GetCommonConfig().RegionName), []*sysRocketMq.Consumer{{
		Topic:   sysRocketMq.AlertContactTopic,
		Handler: consumer.AlertContactHandler,
	}, {
		Topic:   sysRocketMq.RuleTopic,
		Handler: consumer.AlarmRuleHandler,
	}}); err != nil {
		log.Printf("create rocketmq consumer error, %v\n", err)
		return err
	}
	return nil
}

type WebServeLoader struct{}

func (w *WebServeLoader) Load() error {
	return web.Start(config.GetServeConfig())
}

type K8sLoader struct{}

func (k *K8sLoader) Load() error {
	if config.GetCommonConfig().Env != "local" {
		return k8s.InitK8s()
	}
	return nil
}
