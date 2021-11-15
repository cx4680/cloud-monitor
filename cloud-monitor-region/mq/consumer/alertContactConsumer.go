package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func AlertContactConsumer() {
	cfg := config.GetRocketmqConfig()
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(cfg.AlertContactGroup),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.NameServer})),
	)
	var MqMsg forms.MqMsg
	var MsgErr error

	err := c.Subscribe(cfg.AlertContactTopic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
			MsgErr = json.Unmarshal(msgs[i].Body, &MqMsg)
			switch MqMsg.EventEum {
			case enums.InsertAlertContact:
				data, _ := json.Marshal(MqMsg.Data)
				var model models.AlertContact
				MsgErr = json.Unmarshal(data, &model)
				dao.NewAlertContact(database.GetDb()).InsertAlertContact(model)
			case enums.UpdateAlertContact:
				data, _ := json.Marshal(MqMsg.Data)
				var model models.AlertContact
				MsgErr = json.Unmarshal(data, &model)
				dao.NewAlertContact(database.GetDb()).UpdateAlertContact(model)
			case enums.DeleteAlertContact:
				data, _ := json.Marshal(MqMsg.Data)
				var contactId string
				MsgErr = json.Unmarshal(data, &contactId)
				dao.NewAlertContact(database.GetDb()).DeleteAlertContact(contactId)
			case enums.InsertAlertContactInformation:
				data, _ := json.Marshal(MqMsg.Data)
				var model models.AlertContactInformation
				MsgErr = json.Unmarshal(data, &model)
				dao.NewAlertContact(database.GetDb()).InsertAlertContactInformation(model)
			case enums.DeleteAlertContactInformation:
				data, _ := json.Marshal(MqMsg.Data)
				var contactId string
				MsgErr = json.Unmarshal(data, &contactId)
				dao.NewAlertContact(database.GetDb()).DeleteAlertContactInformation(contactId)
			case enums.InsertAlertContactGroupRel:
				data, _ := json.Marshal(MqMsg.Data)
				var model models.AlertContactGroupRel
				MsgErr = json.Unmarshal(data, &model)
				dao.NewAlertContact(database.GetDb()).InsertAlertContactGroupRel(model)
			case enums.DeleteAlertContactGroupRelByContactId:
				data, _ := json.Marshal(MqMsg.Data)
				var contactId string
				MsgErr = json.Unmarshal(data, &contactId)
				dao.NewAlertContact(database.GetDb()).DeleteAlertContactGroupRelByContactId(contactId)
			case enums.CertifyAlertContact:
				data, _ := json.Marshal(MqMsg.Data)
				var activeCode string
				MsgErr = json.Unmarshal(data, &activeCode)
				dao.NewAlertContact(database.GetDb()).CertifyAlertContact(activeCode)
			case enums.InsertAlertContactGroup:
				data, _ := json.Marshal(MqMsg.Data)
				var model models.AlertContactGroup
				MsgErr = json.Unmarshal(data, &model)
				dao.NewAlertContact(database.GetDb()).InsertAlertContactGroup(model)
			case enums.UpdateAlertContactGroup:
				data, _ := json.Marshal(MqMsg.Data)
				var model models.AlertContactGroup
				MsgErr = json.Unmarshal(data, &model)
				dao.NewAlertContact(database.GetDb()).UpdateAlertContactGroup(model)
			case enums.DeleteAlertContactGroup:
				data, _ := json.Marshal(MqMsg.Data)
				var groupId string
				MsgErr = json.Unmarshal(data, &groupId)
				dao.NewAlertContact(database.GetDb()).DeleteAlertContactGroup(groupId)
			case enums.DeleteAlertContactGroupRelByGroupId:
				data, _ := json.Marshal(MqMsg.Data)
				var groupId string
				MsgErr = json.Unmarshal(data, &groupId)
				dao.NewAlertContact(database.GetDb()).DeleteAlertContactGroupRelByGroupId(groupId)
			}
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
