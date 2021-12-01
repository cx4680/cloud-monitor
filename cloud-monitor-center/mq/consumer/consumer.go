package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func InstanceHandler(msgs []*primitive.MessageExt) {
	alarmInstanceDao := dao.AlarmInstance
	for i := range msgs {
		var instances []*models.AlarmInstance
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		err := json.Unmarshal(msgs[i].Body, &instances)
		if err != nil {
			continue
		}
		alarmInstanceDao.UpdateBatchInstanceName(instances)
	}
}

// SmsMarginReminderConsumer 短信余量提醒
func SmsMarginReminderConsumer(msgs []*primitive.MessageExt) {
	svc := service.NewMessageService(messageCenter.NewService())
	for _, msg := range msgs {
		svc.SmsMarginReminder(string(msg.Body))
	}
}

func DeleteInstanceHandler(msgs []*primitive.MessageExt) {
	alarmInstanceDao := dao.AlarmInstance
	for i := range msgs {
		instance := dtos.Instance{}
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		err := json.Unmarshal(msgs[i].Body, &instance)
		if err != nil {
			continue
		}
		alarmInstanceDao.DeleteInstanceList(instance.TenantId, instance.List)
	}
}
