package producer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
)

func SendNotificationRecordMsg(msg []models.NotificationRecord) {
	if msg == nil || len(msg) <= 0 {
		return
	}
	sysRocketMq.SendMsg("notification_sync", tools.ToString(msg))
}

func SendAlertRecordMsg(msg []*models.AlertRecord) {
	sysRocketMq.SendMsg(sysRocketMq.RecordTopic, tools.ToString(msg))
}

func SendInstanceJobMsg(msg []models.AlarmInstance) {
	sysRocketMq.SendMsg(sysRocketMq.InstanceTopic, tools.ToString(msg))
}
