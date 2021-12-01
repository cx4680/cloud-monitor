package loader

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
)

type TransactionLoader struct {
}

func (t *TransactionLoader) Load() error {
	return translate.InitTrans("zh")
}

type TaskLoader struct{}

func (t *TaskLoader) Load() error {
	bt := task.NewBusinessTaskImpl()
	if err := bt.Add(task.BusinessTaskDTO{
		Cron: "0 0 0/1 * * ?",
		Name: "clearAlertRecordJob",
		Task: task.Clear,
	}); err != nil {
		return err
	}

	bt.Start()
	return nil
}

type RocketMQConsumerLoader struct{}

func (r *RocketMQConsumerLoader) Load() error {
	return sysRocketMq.StartConsumersScribe("cloud-monitor-center", []*sysRocketMq.Consumer{{
		Topic:   sysRocketMq.InstanceTopic,
		Handler: consumer.InstanceHandler,
	}, {
		Topic:   sysRocketMq.SmsMarginReminderTopic,
		Handler: consumer.SmsMarginReminderConsumer,
	}, {
		Topic:   sysRocketMq.DeleteInstanceTopic,
		Handler: consumer.DeleteInstanceHandler,
	}})
}

type WebServeLoader struct{}

func (w *WebServeLoader) Load() error {
	return web.Start(config.GetServeConfig())
}
