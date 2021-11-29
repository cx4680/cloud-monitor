package consumer

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonForm "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"

	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
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
			commonDao.AlertContact.Insert(global.DB, model)
		case enums.UpdateAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			commonDao.AlertContact.Update(global.DB, model)
		case enums.DeleteAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model *models.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			commonDao.AlertContact.Delete(global.DB, model)
		case enums.InsertAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var models []*models.AlertContactInformation
			MsgErr = json.Unmarshal(data, &models)
			commonDao.AlertContactInformation.InsertBatch(global.DB, models)
		case enums.DeleteAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var contactId string
			MsgErr = json.Unmarshal(data, &contactId)
			dao.AlertContact.DeleteAlertContactInformation(contactId)
		case enums.InsertAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			//dao.AlertContact.InsertAlertContactGroupRel(model)
		case enums.DeleteAlertContactGroupRelByContactId:
			data, _ := json.Marshal(MqMsg.Data)
			var contactId string
			MsgErr = json.Unmarshal(data, &contactId)
			dao.AlertContact.DeleteAlertContactGroupRelByContactId(contactId)
		case enums.CertifyAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var activeCode string
			MsgErr = json.Unmarshal(data, &activeCode)
			dao.AlertContact.CertifyAlertContact(activeCode)
		case enums.InsertAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			//dao.AlertContact.InsertAlertContactGroup(model)
		case enums.UpdateAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			//dao.AlertContact.UpdateAlertContactGroup(model)
		case enums.DeleteAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var groupId string
			MsgErr = json.Unmarshal(data, &groupId)
			dao.AlertContact.DeleteAlertContactGroup(groupId)
		case enums.DeleteAlertContactGroupRelByGroupId:
			data, _ := json.Marshal(MqMsg.Data)
			var groupId string
			MsgErr = json.Unmarshal(data, &groupId)
			dao.AlertContact.DeleteAlertContactGroupRelByGroupId(groupId)
		}
	}
}
