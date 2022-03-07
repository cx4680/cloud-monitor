package sys_component

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
)

func InitSys() error {

	if err := sys_db.InitDb(config.Cfg.Db); err != nil {
		logger.Logger().Errorf("init database error: %v\n", err)
		return err
	}

	if err := sys_redis.InitClient(config.Cfg.Redis); err != nil {
		logger.Logger().Errorf("init redis error: %v\n", err)
		return err
	}
	if err := sys_rocketmq.InitProducer(); err != nil {
		logger.Logger().Errorf("create rocketmq consumer error, %v\n", err)
		return err
	}

	topics := []sys_rocketmq.Topic{
		sys_rocketmq.SmsMarginReminderTopic,
		sys_rocketmq.InstanceTopic,
		sys_rocketmq.ContactTopic,
		sys_rocketmq.ContactGroupTopic,
		sys_rocketmq.RecordTopic,
		sys_rocketmq.AlarmInfoTopic,
		sys_rocketmq.RuleTopic,
		sys_rocketmq.NotificationSyncTopic,
		sys_rocketmq.DeleteInstanceTopic,
		sys_rocketmq.MonitorProductTopic,
		sys_rocketmq.MonitorItemTopic,
	}

	err := sys_rocketmq.CreateTopics(topics...)
	if err != nil {
		logger.Logger().Errorf("create rocketmq topics error, %v\n", err)
		return err
	}
	return nil
}
