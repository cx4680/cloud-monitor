package sysComponent

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRedis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"log"
)

func InitSys() error {

	if err := sysDb.InitDb(config.GetDbConfig()); err != nil {
		log.Printf("init database error: %v\n", err)
		return err
	}

	if err := sysRedis.InitClient(config.GetRedisConfig()); err != nil {
		log.Printf("init redis error: %v\n", err)
		return err
	}
	if err := sysRocketMq.InitProducer(); err != nil {
		log.Printf("create rocketmq consumer error, %v\n", err)
		return err
	}

	topics := []sysRocketMq.Topic{
		sysRocketMq.SmsMarginReminderTopic,
		sysRocketMq.InstanceTopic,
		sysRocketMq.AlertContactTopic,
		sysRocketMq.RecordTopic,
		sysRocketMq.RuleTopic,
		sysRocketMq.NotificationSyncTopic,
		sysRocketMq.DeleteInstanceTopic,
	}

	err := sysRocketMq.CreateTopics(topics...)
	if err != nil {
		log.Printf("create rocketmq topics error, %v\n", err)
		return err
	}
	return nil
}
