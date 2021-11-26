package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

//TODO 需要初始化加载
// SmsMarginReminderConsumer 短信余量提醒
func SmsMarginReminderConsumer(svc *service.MessageService, msgs []*primitive.MessageExt) {
	for _, msg := range msgs {
		svc.SmsMarginReminder(string(msg.Body))
	}
}
