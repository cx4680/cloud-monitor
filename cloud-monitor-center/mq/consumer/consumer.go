package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func InstanceHandler(msgs []*primitive.MessageExt) {
	alarmInstanceDao := dao.AlarmInstance
	for i := range msgs {
		var instances []*model.AlarmInstance
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
	svc := service.NewMessageService(message_center.NewService())
	for _, msg := range msgs {
		svc.SmsMarginReminder(string(msg.Body))
	}
}

func DeleteInstanceHandler(msgs []*primitive.MessageExt) {
	alarmInstanceDao := dao.AlarmInstance
	for i := range msgs {
		instance := dto.Instance{}
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		err := json.Unmarshal(msgs[i].Body, &instance)
		if err != nil {
			continue
		}
		alarmInstanceDao.DeleteInstanceList(instance.TenantId, instance.List)
	}
}

func AlertRecordAddHandler(msgs []*primitive.MessageExt) {
	for _, msg := range msgs {
		var list []model.AlertRecord
		jsonutil.ToObject(string(msg.Body), &list)
		if list != nil && len(list) > 0 {
			dao.AlertRecord.InsertBatch(global.DB, list)
		}
	}
}
