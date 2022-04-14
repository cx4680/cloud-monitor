package producer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/publisher"
)

func SendInstanceJobMsg(topic sys_rocketmq.Topic, msg interface{}) {
	publisher.GlobalPublisher.Pub(publisher.PubMessage{
		Topic: topic,
		Data:  jsonutil.ToString(msg),
	})
}
