package consumer

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonForm "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
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
			var model models.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.NewAlertContact().InsertAlertContact(model)
		case enums.UpdateAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContact
			MsgErr = json.Unmarshal(data, &model)
			dao.NewAlertContact().UpdateAlertContact(model)
		case enums.DeleteAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var contactId string
			MsgErr = json.Unmarshal(data, &contactId)
			dao.NewAlertContact().DeleteAlertContact(contactId)
		case enums.InsertAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContactInformation
			MsgErr = json.Unmarshal(data, &model)
			dao.NewAlertContact().InsertAlertContactInformation(model)
		case enums.DeleteAlertContactInformation:
			data, _ := json.Marshal(MqMsg.Data)
			var contactId string
			MsgErr = json.Unmarshal(data, &contactId)
			dao.NewAlertContact().DeleteAlertContactInformation(contactId)
		case enums.InsertAlertContactGroupRel:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContactGroupRel
			MsgErr = json.Unmarshal(data, &model)
			dao.NewAlertContact().InsertAlertContactGroupRel(model)
		case enums.DeleteAlertContactGroupRelByContactId:
			data, _ := json.Marshal(MqMsg.Data)
			var contactId string
			MsgErr = json.Unmarshal(data, &contactId)
			dao.NewAlertContact().DeleteAlertContactGroupRelByContactId(contactId)
		case enums.CertifyAlertContact:
			data, _ := json.Marshal(MqMsg.Data)
			var activeCode string
			MsgErr = json.Unmarshal(data, &activeCode)
			dao.NewAlertContact().CertifyAlertContact(activeCode)
		case enums.InsertAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.NewAlertContact().InsertAlertContactGroup(model)
		case enums.UpdateAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var model models.AlertContactGroup
			MsgErr = json.Unmarshal(data, &model)
			dao.NewAlertContact().UpdateAlertContactGroup(model)
		case enums.DeleteAlertContactGroup:
			data, _ := json.Marshal(MqMsg.Data)
			var groupId string
			MsgErr = json.Unmarshal(data, &groupId)
			dao.NewAlertContact().DeleteAlertContactGroup(groupId)
		case enums.DeleteAlertContactGroupRelByGroupId:
			data, _ := json.Marshal(MqMsg.Data)
			var groupId string
			MsgErr = json.Unmarshal(data, &groupId)
			dao.NewAlertContact().DeleteAlertContactGroupRelByGroupId(groupId)
		}
	}
}

func AlarmRuleHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		json.Unmarshal(msgs[i].Body, &MqMsg)
		data, _ := json.Marshal(MqMsg.Data)
		var tenantId string
		tx(MqMsg, data, tenantId)
	}
}

func tx(MqMsg forms.MqMsg, data []byte, tenantId string) error {
	tx := database.GetDb()
	var err error
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Errorf("%v", err)
			tx.Rollback()
			err = fmt.Errorf("%v", err)
		}
	}()
	handleMsg(MqMsg, data, tenantId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}

func handleMsg(MqMsg forms.MqMsg, data []byte, tenantId string) {
	ruleDao := commonDao.NewAlarmRuleDao()
	instanceDao := commonDao.NewInstanceDao()
	prometheusDao := dao.NewPrometheusRuleDao()

	switch MqMsg.EventEum {
	case enums.CreateRule:
		var param commonForm.AlarmRuleAddReqDTO
		json.Unmarshal(data, &param)
		ruleDao.SaveRule(database.GetDb(), &param)
		tenantId = param.TenantId
	case enums.UpdateRule:
		var param commonForm.AlarmRuleAddReqDTO
		json.Unmarshal(data, &param)
		ruleDao.UpdateRule(database.GetDb(), &param)
		tenantId = param.TenantId
	case enums.EnableRule:
		var param commonForm.RuleReqDTO
		json.Unmarshal(data, &param)
		ruleDao.UpdateRuleState(database.GetDb(), &param)
	case enums.DisableRule:
		var param commonForm.RuleReqDTO
		json.Unmarshal(data, &param)
		ruleDao.UpdateRuleState(database.GetDb(), &param)
	case enums.DeleteRule:
		var param commonForm.RuleReqDTO
		json.Unmarshal(data, &param)
		ruleDao.DeleteRule(database.GetDb(), &param)
	case enums.UnbindRule:
		var param commonForm.UnBindRuleParam
		json.Unmarshal(data, &param)
		instanceDao.UnbindInstance(database.GetDb(), &param)
	case enums.BindRule:
		var param commonForm.InstanceBindRuleDTO
		json.Unmarshal(data, &param)
		instanceDao.BindInstance(database.GetDb(), &param)
	default:
		logger.Logger().Warnf("不支持的消息类型，消息类型：%v,消息%s", MqMsg.EventEum, string(data))
	}
	if len(tenantId) > 0 {
		prometheusDao.GenerateUserPrometheusRule("", "", tenantId)
	}
}
