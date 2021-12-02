package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
)

type TransactionActuatorStage struct {
}

func (s *TransactionActuatorStage) Exec(c *context.Context) error {
	return translate.InitTrans("zh")
}

type TaskActuatorStage struct {
}

func (ta *TaskActuatorStage) Exec(c *context.Context) error {
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

type MQActuatorStage struct {
}

func (ma *MQActuatorStage) Exec(c *context.Context) error {
	return sysRocketMq.StartConsumersScribe(sysRocketMq.Group(config.GetCommonConfig().RegionName), []*sysRocketMq.Consumer{{
		Topic:   sysRocketMq.AlertContactTopic,
		Handler: consumer.AlertContactHandler,
	}, {
		Topic:   sysRocketMq.RuleTopic,
		Handler: consumer.AlarmRuleHandler,
	}})
}

type K8sActuatorStage struct{}

func (k *K8sActuatorStage) Exec(c *context.Context) error {
	return k8s.InitK8s()
}

type WebActuatorStage struct {
}

func (wa *WebActuatorStage) Exec(c *context.Context) error {
	return web.Start(config.GetServeConfig())
}
