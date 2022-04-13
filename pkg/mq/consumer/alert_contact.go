package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
)

func ContactHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		var MsgErr error
		var MqMsg form.MqMsg
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		MsgErr = json.Unmarshal(msgs[i].Body, &MqMsg)
		switch MqMsg.EventEum {
		case enum.InsertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var contactMsg *service.ContactMsg
			MsgErr = json.Unmarshal(data, &contactMsg)
			MsgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Insert(db, contactMsg.Contact)
				dao.ContactInformation.InsertBatch(db, contactMsg.ContactInformationList)
				dao.ContactGroupRel.InsertBatch(db, contactMsg.ContactGroupRelList)
				return nil
			})
		case enum.UpdateContact:
			data, _ := json.Marshal(MqMsg.Data)
			var contactMsg *service.ContactMsg
			MsgErr = json.Unmarshal(data, &contactMsg)
			MsgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Update(db, contactMsg.Contact)
				dao.ContactInformation.Update(db, contactMsg.ContactInformationList)
				dao.ContactGroupRel.Update(db, contactMsg.ContactGroupRelList, contactMsg.Param)
				return nil
			})
		case enum.DeleteContact:
			data, _ := json.Marshal(MqMsg.Data)
			var contactMsg *service.ContactMsg
			MsgErr = json.Unmarshal(data, &contactMsg)
			MsgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Delete(db, contactMsg.Contact)
				dao.ContactInformation.Delete(db, contactMsg.ContactInformation)
				dao.ContactGroupRel.Delete(db, contactMsg.ContactGroupRel)
				return nil
			})
		case enum.ActivateContact:
			data, _ := json.Marshal(MqMsg.Data)
			var activeCode string
			MsgErr = json.Unmarshal(data, &activeCode)
			dao.Contact.ActivateContact(global.DB, activeCode)
		}
		if MsgErr != nil {
			logger.Logger().Errorf("%v", MsgErr)
		}
	}
}

func ContactGroupHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		var MsgErr error
		var MqMsg form.MqMsg
		MsgErr = json.Unmarshal(msgs[i].Body, &MqMsg)
		switch MqMsg.EventEum {
		case enum.InsertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var contactGroupMsg *service.ContactGroupMsg
			MsgErr = json.Unmarshal(data, &contactGroupMsg)
			MsgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Insert(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.InsertBatch(db, contactGroupMsg.ContactGroupRelList)
				return nil
			})
		case enum.UpdateContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var contactGroupMsg *service.ContactGroupMsg
			MsgErr = json.Unmarshal(data, &contactGroupMsg)
			MsgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Update(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.Update(db, contactGroupMsg.ContactGroupRelList, contactGroupMsg.Param)
				return nil
			})
		case enum.DeleteContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var contactGroupMsg *service.ContactGroupMsg
			MsgErr = json.Unmarshal(data, &contactGroupMsg)
			MsgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Delete(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.Delete(db, contactGroupMsg.ContactGroupRel)
				return nil
			})
		}
		if MsgErr != nil {
			logger.Logger().Errorf("%v", MsgErr)
		}
	}
}
