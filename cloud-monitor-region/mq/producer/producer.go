package producer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
)

func SendNotificationRecordMsg(msg []models.NotificationRecord) {
	if msg == nil || len(msg) <= 0 {
		return
	}
	mq.SendMsg("notification_sync", tools.ToString(msg))
}

func SendAlertRecordMsg(msg []*models.AlertRecord) {
	cfg := config.GetRocketmqConfig()
	mq.SendMsg(cfg.RecordTopic, tools.ToString(msg))
}
