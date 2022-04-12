package subscriber

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq/consumer"
)

type Subscriber struct {
}

func (*Subscriber) Run() error {
	return sys_rocketmq.StartConsumersScribe(sys_rocketmq.Group(config.Cfg.Common.RegionName),
		[]*sys_rocketmq.Consumer{{
			Topic:   sys_rocketmq.InstanceTopic,
			Handler: consumer.InstanceHandler,
		}, {
			Topic:   sys_rocketmq.SmsMarginReminderTopic,
			Handler: consumer.SmsMarginReminderConsumer,
		}, {
			Topic:   sys_rocketmq.DeleteInstanceTopic,
			Handler: consumer.DeleteInstanceHandler,
		}, {
			Topic:   sys_rocketmq.AlarmTopic,
			Handler: consumer.AlarmAddHandler,
		}, {
			Topic:   sys_rocketmq.ContactTopic,
			Handler: consumer.ContactHandler,
		}, {
			Topic:   sys_rocketmq.ContactGroupTopic,
			Handler: consumer.ContactGroupHandler,
		}, {
			Topic:   sys_rocketmq.RuleTopic,
			Handler: consumer.AlarmRuleHandler,
		}, {
			Topic:   sys_rocketmq.MonitorProductTopic,
			Handler: consumer.MonitorProductHandler,
		}, {
			Topic:   sys_rocketmq.MonitorItemTopic,
			Handler: consumer.MonitorItemHandler,
		}})

}
