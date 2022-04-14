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
		var msgErr error
		var mqMsg form.MqMsg
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		msgErr = json.Unmarshal(msgs[i].Body, &mqMsg)
		switch mqMsg.EventEum {
		case enum.InsertContact:
			contactMsg := buildContactData(mqMsg.Data)
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Insert(db, contactMsg.Contact)
				dao.ContactInformation.InsertBatch(db, contactMsg.ContactInformationList)
				dao.ContactGroupRel.InsertBatch(db, contactMsg.ContactGroupRelList)
				return nil
			})
		case enum.UpdateContact:
			contactMsg := buildContactData(mqMsg.Data)
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Update(db, contactMsg.Contact)
				dao.ContactInformation.Update(db, contactMsg.ContactInformationList)
				dao.ContactGroupRel.Update(db, contactMsg.ContactGroupRelList, contactMsg.Param)
				return nil
			})
		case enum.DeleteContact:
			contactMsg := buildContactData(mqMsg.Data)
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Delete(db, contactMsg.Contact)
				dao.ContactInformation.Delete(db, contactMsg.ContactInformation)
				dao.ContactGroupRel.Delete(db, contactMsg.ContactGroupRel)
				return nil
			})
		case enum.ActivateContact:
			data, _ := json.Marshal(mqMsg.Data)
			var activeCode string
			msgErr = json.Unmarshal(data, &activeCode)
			dao.Contact.ActivateContact(global.DB, activeCode)
		}
		if msgErr != nil {
			logger.Logger().Errorf("%v", msgErr)
		}
	}
}

func ContactGroupHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		var msgErr error
		var mqMsg form.MqMsg
		msgErr = json.Unmarshal(msgs[i].Body, &mqMsg)
		switch mqMsg.EventEum {
		case enum.InsertContactGroup:
			contactGroupMsg := buildContactGroupData(mqMsg.Data)
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Insert(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.InsertBatch(db, contactGroupMsg.ContactGroupRelList)
				return nil
			})
		case enum.UpdateContactGroup:
			contactGroupMsg := buildContactGroupData(mqMsg.Data)
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Update(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.Update(db, contactGroupMsg.ContactGroupRelList, contactGroupMsg.Param)
				return nil
			})
		case enum.DeleteContactGroup:
			contactGroupMsg := buildContactGroupData(mqMsg.Data)
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Delete(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.Delete(db, contactGroupMsg.ContactGroupRel)
				return nil
			})
		}
		if msgErr != nil {
			logger.Logger().Errorf("%v", msgErr)
		}
	}
}

func buildContactData(msgData interface{}) *service.ContactMsg {
	data, _ := json.Marshal(msgData)
	var contactMsg *service.ContactMsg
	err := json.Unmarshal(data, &contactMsg)
	if err != nil {
		logger.Logger().Error(err)
		return nil
	}
	contactMsg.Contact.Id = 0
	for i := range contactMsg.ContactInformationList {
		contactMsg.ContactInformationList[i].Id = 0
	}
	for i := range contactMsg.ContactGroupRelList {
		contactMsg.ContactGroupRelList[i].Id = 0
	}
	return contactMsg
}

func buildContactGroupData(msgData interface{}) *service.ContactGroupMsg {
	data, _ := json.Marshal(msgData)
	var contactGroupMsg *service.ContactGroupMsg
	err := json.Unmarshal(data, &contactGroupMsg)
	if err != nil {
		logger.Logger().Error(err)
		return nil
	}
	contactGroupMsg.ContactGroup.Id = 0
	for _, v := range contactGroupMsg.ContactGroupRelList {
		v.Id = 0
	}
	return contactGroupMsg
}
