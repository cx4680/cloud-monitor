package sysComponent

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRedis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
)

func InitSys() error {

	if err := sysDb.InitDb(config.GetDbConfig()); err != nil {
		logger.Logger().Errorf("init database error: %v\n", err)
		return err
	}

	if err := sysRedis.InitClient(config.GetRedisConfig()); err != nil {
		logger.Logger().Errorf("init redis error: %v\n", err)
		return err
	}
	if err := sysRocketMq.InitProducer(); err != nil {
		logger.Logger().Errorf("create rocketmq consumer error, %v\n", err)
		return err
	}

	topics := []sysRocketMq.Topic{
		sysRocketMq.SmsMarginReminderTopic,
		sysRocketMq.InstanceTopic,
		sysRocketMq.AlertContactTopic,
		sysRocketMq.AlertContactGroupTopic,
		sysRocketMq.RecordTopic,
		sysRocketMq.RuleTopic,
		sysRocketMq.NotificationSyncTopic,
		sysRocketMq.DeleteInstanceTopic,
	}

	err := sysRocketMq.CreateTopics(topics...)
	if err != nil {
		logger.Logger().Errorf("create rocketmq topics error, %v\n", err)
		return err
	}
	return nil
}
