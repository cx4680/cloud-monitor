package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var MqMsg forms.MqMsg
var MsgErr error

func AlertContactHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		MsgErr = json.Unmarshal(msgs[i].Body, &MqMsg)
		switch MqMsg.EventEum {
		case enums.InsertAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContact.Insert(global.DB, model)
		case enums.UpdateAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContact.Update(global.DB, model)
		case enums.DeleteAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContact.Delete(global.DB, model)
		case enums.InsertAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*models.AlertContactInformation
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactInformation.InsertBatch(global.DB, modelList)
		case enums.UpdateAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*models.AlertContactInformation
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactInformation.Update(global.DB, modelList)
		case enums.DeleteAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContactInformation
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactInformation.Delete(global.DB, model)
		case enums.InsertAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*models.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactGroupRel.InsertBatch(global.DB, modelList)
		case enums.UpdateAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var modelList []*models.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &modelList)
			dao.AlertContactGroupRel.Update(global.DB, modelList)
		case enums.DeleteAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroupRel.Delete(global.DB, model)
		case enums.CertifyAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var activeCode string
			MsgErr = json.Unmarshal(data, &activeCode)
			dao.AlertContact.CertifyAlertContact(activeCode)
		case enums.InsertAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroup.Insert(global.DB, model)
		case enums.UpdateAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroup.Update(global.DB, model)
		case enums.DeleteAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.AlertContactGroup.Delete(global.DB, model)
		}
	}
}
