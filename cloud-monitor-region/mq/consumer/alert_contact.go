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

func AlertContactHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		var MsgErr error
		var MqMsg form.MqMsg
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		MsgErr = json.Unmarshal(msgs[i].Body, &MqMsg)
		switch MqMsg.EventEum {
		case enum.InsertAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContact.Insert(global.DB, model)
		case enum.UpdateAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContact.Update(global.DB, model)
		case enum.DeleteAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContact.Delete(global.DB, model)
		case enum.InsertAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.AlertContactInformation
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactInformation.InsertBatch(global.DB, modelList)
		case enum.UpdateAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.AlertContactInformation
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactInformation.Update(global.DB, modelList)
		case enum.DeleteAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContactInformation
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactInformation.Delete(global.DB, model)
		case enum.CertifyAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var activeCode string
			MsgErr = json.Unmarshal(data, &activeCode)
			dao.AlertContact.CertifyAlertContact(global.DB, activeCode)
		case enum.InsertAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactGroupRel.InsertBatch(global.DB, modelList)
		case enum.UpdateAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model model.UpdateAlertContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroupRel.Update(global.DB, model.RelList, model.Param)
		case enum.DeleteAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroupRel.Delete(global.DB, model)
		}
		if MsgErr != nil {
			logger.Logger().Errorf("%v", MsgErr)
		}
	}
}

func AlertContactGroupHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		var MsgErr error
		var MqMsg form.MqMsg
		MsgErr = json.Unmarshal(msgs[i].Body, &MqMsg)
		switch MqMsg.EventEum {
		case enum.InsertAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*model.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactGroupRel.InsertBatch(global.DB, modelList)
		case enum.UpdateAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model model.UpdateAlertContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroupRel.Update(global.DB, model.RelList, model.Param)
		case enum.DeleteAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroupRel.Delete(global.DB, model)
		case enum.InsertAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroup.Insert(global.DB, model)
		case enum.UpdateAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroup.Update(global.DB, model)
		case enum.DeleteAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *model.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroup.Delete(global.DB, model)
		}
		if MsgErr != nil {
			logger.Logger().Errorf("%v", MsgErr)
		}
	}
}
