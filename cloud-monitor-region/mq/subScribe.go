package mq

import "code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"

func SubScribe() {
	consumer.AlertContactConsumer()
	consumer.AlarmRuleConsumer()
}
