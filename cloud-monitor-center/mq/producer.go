package mq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
)

func SendMsg(topic sysRocketMq.Topic, eventEum enums.EventEum, module interface{}) error {
	var mqMsg = forms.MqMsg{
		EventEum: eventEum,
		Data:     module,
	}
	str := tools.ToString(mqMsg)
	return sysRocketMq.SendMsg(topic, str)
}
