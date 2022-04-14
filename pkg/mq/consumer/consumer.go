package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	dto2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/message_center"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
)

func InstanceHandler(msgs []*primitive.MessageExt) {
	alarmInstanceDao := dao2.AlarmInstance
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
	alarmInstanceDao := dao2.AlarmInstance
	for i := range msgs {
		instance := dto2.Instance{}
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		err := json.Unmarshal(msgs[i].Body, &instance)
		if err != nil {
			continue
		}
		alarmInstanceDao.DeleteInstanceList(instance.TenantId, instance.List)
	}
}

func AlarmAddHandler(msgs []*primitive.MessageExt) {
	for _, msg := range msgs {
		var data dto2.AlarmSyncData
		if err := jsonutil.ToObjectWithError(string(msg.Body), &data); err != nil {
			logger.Logger().Errorf("序列化数据失败, %v", msg)
		}
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			if len(data.RecordList) > 0 {
				for i, _ := range data.RecordList {
					data.RecordList[i].Id = 0
				}
				dao2.AlarmRecord.InsertBatch(tx, data.RecordList)
			}
			if len(data.InfoList) > 0 {
				for i, _ := range data.InfoList {
					data.InfoList[i].Id = 0
				}
				dao2.AlarmInfo.InsertBatch(tx, data.InfoList)
			}
			return nil
		}); err != nil {
			logger.Logger().Errorf("消费MQ数据失败, %v", msg)
		}

	}

}
