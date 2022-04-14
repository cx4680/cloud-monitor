package mq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/publisher"
)

func SendMsg(topic sys_rocketmq.Topic, eventEum enum.EventEum, module interface{}) error {
	var mqMsg = form.MqMsg{
		EventEum: eventEum,
		Data:     module,
	}
	str := jsonutil.ToString(mqMsg)
	return publisher.GlobalPublisher.Pub(publisher.PubMessage{
		Topic: topic,
		Data:  str,
	})
}
