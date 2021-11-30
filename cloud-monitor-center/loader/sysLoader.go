package loader

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/loader"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
)

type SysLoaderImpl struct {
	loader.AbstractSysLoader
}

func (s SysLoaderImpl) InitTask() error {
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

func (s *SysLoaderImpl) InitRocketMqConsumers() error {
	//TODO 初始化消费者
	if err := sysRocketMq.StartConsumersScribe("cloud-monitor-center", []*sysRocketMq.Consumer{{
		Topic:   sysRocketMq.InstanceTopic,
		Handler: consumer.InstanceHandler,
	}}); err != nil {
		return err
	}
	return nil
}

func (s *SysLoaderImpl) InitWebServe() error {
	if err := web.Start(config.GetServeConfig()); err != nil {
		return err
	}
	return nil
}

func (s *SysLoaderImpl) InitTrans() error {
	return translate.InitTrans("zh")
}
