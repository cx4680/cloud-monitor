package loader

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/loader"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"log"
)

type SysLoaderImpl struct {
	loader.AbstractSysLoader
}

func (s *SysLoaderImpl) InitTask() error {
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

func (s *SysLoaderImpl) InitRocketMqConsumers() error {
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

func (s *SysLoaderImpl) InitWebServe() error {
	if err := web.Start(config.GetServeConfig()); err != nil {
		return err
	}
	return nil
}

func (s *SysLoaderImpl) InitTrans() error {
	return translate.InitTrans("zh")
}
