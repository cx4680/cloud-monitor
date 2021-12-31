package producer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
)

func SendNotificationRecordMsg(msg []model.NotificationRecord) {
	if msg == nil || len(msg) <= 0 {
		return
	}
	sys_rocketmq.SendMsg("notification_sync", jsonutil.ToString(msg))
}

func SendAlertRecordMsg(msg []*model.AlertRecord) {
	sys_rocketmq.SendMsg(sys_rocketmq.RecordTopic, jsonutil.ToString(msg))
}

func SendInstanceJobMsg(topic sys_rocketmq.Topic, msg interface{}) {
	sys_rocketmq.SendMsg(topic, jsonutil.ToString(msg))
}
