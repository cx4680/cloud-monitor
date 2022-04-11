package mq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
)

func SendMsg(topic sys_rocketmq.Topic, eventEum enum.EventEum, module interface{}) error {
	var mqMsg = form.MqMsg{
		EventEum: eventEum,
		Data:     module,
	}
	str := jsonutil.ToString(mqMsg)
	return sys_rocketmq.SendMsg(topic, str)
}
