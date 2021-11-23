package producer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
)

func SendNotificationRecordMsg(msg []models.NotificationRecord) {
	if msg == nil || len(msg) <= 0 {
		return
	}
	sysRocketMq.SendMsg("notification_sync", tools.ToString(msg))
}

func SendAlertRecordMsg(msg []*models.AlertRecord) {
	cfg := config.GetRocketmqConfig()
	sysRocketMq.SendMsg(cfg.RecordTopic, tools.ToString(msg))
}

func SendInstanceJobMsg(msg []models.AlarmInstance) {
	cfg := config.GetRocketmqConfig()
	sysRocketMq.SendMsg(cfg.InstanceTopic, tools.ToString(msg))
}
