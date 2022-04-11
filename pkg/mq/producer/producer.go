package producer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
)

func SendInstanceJobMsg(topic sys_rocketmq.Topic, msg interface{}) {
	sys_rocketmq.SendMsg(topic, jsonutil.ToString(msg))
}
