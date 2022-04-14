package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
)

func ContactHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		var msgErr error
		var mqMsg form.MqMsg
		logger.Logger().Infof("subscribe callback: %v \n", msgs[i])
		msgErr = jsonutil.ToObjectWithError(string(msgs[i].Body), &mqMsg)
		if msgErr != nil {
			continue
		}
		switch mqMsg.EventEum {
		case enum.InsertContact:
			contactMsg, err := buildContactData(mqMsg.Data)
			if err != nil {
				continue
			}
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Insert(db, contactMsg.Contact)
				dao.ContactInformation.InsertBatch(db, contactMsg.ContactInformationList)
				dao.ContactGroupRel.InsertBatch(db, contactMsg.ContactGroupRelList)
				return nil
			})
		case enum.UpdateContact:
			contactMsg, err := buildContactData(mqMsg.Data)
			if err != nil {
				continue
			}
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Update(db, contactMsg.Contact)
				dao.ContactInformation.Update(db, contactMsg.ContactInformationList)
				dao.ContactGroupRel.Update(db, contactMsg.ContactGroupRelList, contactMsg.Param)
				return nil
			})
		case enum.DeleteContact:
			contactMsg, err := buildContactData(mqMsg.Data)
			if err != nil {
				continue
			}
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.Contact.Delete(db, contactMsg.Contact)
				dao.ContactInformation.Delete(db, contactMsg.ContactInformation)
				dao.ContactGroupRel.Delete(db, contactMsg.ContactGroupRel)
				return nil
			})
		case enum.ActivateContact:
			data := jsonutil.ToString(mqMsg.Data)
			var activeCode string
			msgErr = jsonutil.ToObjectWithError(data, &activeCode)
			dao.Contact.ActivateContact(global.DB, activeCode)
		}
		if msgErr != nil {
			logger.Logger().Errorf("%v", msgErr)
		}
	}
}

func ContactGroupHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		logger.Logger().Infof("subscribe callback: %v \n", msgs[i])
		var msgErr error
		var mqMsg form.MqMsg
		msgErr = jsonutil.ToObjectWithError(string(msgs[i].Body), &mqMsg)
		if msgErr != nil {
			continue
		}
		switch mqMsg.EventEum {
		case enum.InsertContactGroup:
			contactGroupMsg, err := buildContactGroupData(mqMsg.Data)
			if err != nil {
				continue
			}
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Insert(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.InsertBatch(db, contactGroupMsg.ContactGroupRelList)
				return nil
			})
		case enum.UpdateContactGroup:
			contactGroupMsg, err := buildContactGroupData(mqMsg.Data)
			if err != nil {
				continue
			}
			msgErr = global.DB.Transaction(func(db *gorm.DB) error {
				dao.ContactGroup.Update(db, contactGroupMsg.ContactGroup)
				dao.ContactGroupRel.Update(db, contactGroupMsg.ContactGroupRelList, contactGroupMsg.Param)
				return nil
			})
		case enum.DeleteContactGroup:
			contactGroupMsg, err := buildContactGroupData(mqMsg.Data)
			if err != nil {
				continue
			}
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

func buildContactData(msgData interface{}) (*service.ContactMsg, error) {
	data := jsonutil.ToString(msgData)
	var contactMsg *service.ContactMsg
	err := jsonutil.ToObjectWithError(data, &contactMsg)
	if err != nil {
		return nil, err
	}
	contactMsg.Contact.Id = 0
	for i := range contactMsg.ContactInformationList {
		contactMsg.ContactInformationList[i].Id = 0
	}
	for i := range contactMsg.ContactGroupRelList {
		contactMsg.ContactGroupRelList[i].Id = 0
	}
	return contactMsg, nil
}

func buildContactGroupData(msgData interface{}) (*service.ContactGroupMsg, error) {
	data := jsonutil.ToString(msgData)
	var contactGroupMsg *service.ContactGroupMsg
	err := jsonutil.ToObjectWithError(data, &contactGroupMsg)
	if err != nil {
		return nil, err
	}
	contactGroupMsg.ContactGroup.Id = 0
	for _, v := range contactGroupMsg.ContactGroupRelList {
		v.Id = 0
	}
	return contactGroupMsg, nil
}
