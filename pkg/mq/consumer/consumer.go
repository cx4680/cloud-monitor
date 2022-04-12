package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
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

func AlarmAddHandler(msgs []*primitive.MessageExt) {
	for _, msg := range msgs {
		var data dto.AlarmSyncData
		if err := jsonutil.ToObjectWithError(string(msg.Body), &data); err != nil {
			logger.Logger().Errorf("序列化数据失败, %v", msg)
		}
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			if len(data.RecordList) > 0 {
				dao.AlarmRecord.InsertBatch(tx, data.RecordList)
			}
			if len(data.InfoList) > 0 {
				dao.AlarmInfo.InsertBatch(tx, data.InfoList)
			}
			return nil
		}); err != nil {
			logger.Logger().Errorf("消费MQ数据失败, %v", msg)
		}

	}

}
