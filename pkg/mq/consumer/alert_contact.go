package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
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
			var model *model.Contact
			MsgErr = json.Unmarshal(data, &model)
			dao.Contact.Insert(global.DB, model)
		case enum.UpdateContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.Contact
			MsgErr = json.Unmarshal(data, &model)
			dao.Contact.Update(global.DB, model)
		case enum.DeleteContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.Contact
			MsgErr = json.Unmarshal(data, &model)
			dao.Contact.Delete(global.DB, model)
		case enum.InsertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.ContactInformation
			MsgErr = json.Unmarshal(data, &modelList)
			dao.ContactInformation.InsertBatch(global.DB, modelList)
		case enum.UpdateContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.ContactInformation
			MsgErr = json.Unmarshal(data, &modelList)
			dao.ContactInformation.Update(global.DB, modelList)
		case enum.DeleteContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.ContactInformation
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactInformation.Delete(global.DB, model)
		case enum.ActivateContact:
			data, _ := json.Marshal(MqMsg.Data)
			var activeCode string
			MsgErr = json.Unmarshal(data, &activeCode)
			dao.Contact.ActivateContact(global.DB, activeCode)
		case enum.InsertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.ContactGroupRel
			MsgErr = json.Unmarshal(data, &modelList)
			dao.ContactGroupRel.InsertBatch(global.DB, modelList)
		case enum.UpdateContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model model.UpdateContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactGroupRel.Update(global.DB, model.RelList, model.Param)
		case enum.DeleteContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.ContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactGroupRel.Delete(global.DB, model)
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
		case enum.InsertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.ContactGroupRel
			MsgErr = json.Unmarshal(data, &modelList)
			dao.ContactGroupRel.InsertBatch(global.DB, modelList)
		case enum.UpdateContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model model.UpdateContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactGroupRel.Update(global.DB, model.RelList, model.Param)
		case enum.DeleteContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.ContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactGroupRel.Delete(global.DB, model)
		case enum.InsertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.ContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactGroup.Insert(global.DB, model)
		case enum.UpdateContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.ContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactGroup.Update(global.DB, model)
		case enum.DeleteContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.ContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.ContactGroup.Delete(global.DB, model)
		}
		if MsgErr != nil {
			logger.Logger().Errorf("%v", MsgErr)
		}
	}
}
