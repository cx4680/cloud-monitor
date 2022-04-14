package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func AlarmRuleHandler(msgList []*primitive.MessageExt) {
	for i := range msgList {
		logger.Logger().Infof("subscribe callback: %v \n", msgList[i])
		var MqMsg form.MqMsg
		if err := json.Unmarshal(msgList[i].Body, &MqMsg); err != nil {
			continue
		}
		data, err := json.Marshal(MqMsg.Data)
		if err != nil {
			continue
		}
		handleMsg(MqMsg, data)
	}
}

func handleMsg(MqMsg form.MqMsg, data []byte) {
	ruleDao := dao.AlarmRule
	instanceDao := dao.Instance
	prometheusDao := service.PrometheusRule
	var tenantId string
	switch MqMsg.EventEum {
	case enum.CreateRule:
		var param form.AlarmRuleAddReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.SaveRule(global.DB, &param)
		tenantId = param.TenantId
	case enum.UpdateRule:
		var param form.AlarmRuleAddReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.UpdateRule(global.DB, &param)
		tenantId = param.TenantId
	case enum.ChangeStatus:
		var param form.RuleReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.UpdateRuleState(global.DB, &param)
		tenantId = param.TenantId
	case enum.DeleteRule:
		var param form.RuleReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.DeleteRule(global.DB, &param)
		tenantId = param.TenantId
	case enum.UnbindRule:
		var param form.UnBindRuleParam
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		instanceDao.UnbindInstance(global.DB, &param)
		tenantId = param.TenantId
	case enum.BindRule:
		var param form.InstanceBindRuleDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		instanceDao.BindInstance(global.DB, &param)
		tenantId = param.TenantId
	default:
		logger.Logger().Warnf("不支持的消息类型，消息类型：%v,消息%s", MqMsg.EventEum, string(data))
	}
	if len(tenantId) > 0 {
		prometheusDao.GenerateUserPrometheusRule(tenantId)
	}
}
