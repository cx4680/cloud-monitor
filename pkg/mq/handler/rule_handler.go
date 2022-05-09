package handler

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
	"encoding/json"
)

func HandleMsg(event enum.EventEum, data []byte, needSave bool) {
	ruleDao := dao.AlarmRule
	instanceDao := dao.Instance
	prometheusDao := k8s.PrometheusRule
	var tenantId string
	switch event {
	case enum.CreateRule:
		var param form.AlarmRuleAddReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		if needSave {
			ruleDao.SaveRule(global.DB, &param)
		}
		tenantId = param.TenantId
	case enum.UpdateRule:
		var param form.AlarmRuleAddReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		if needSave {
			ruleDao.UpdateRule(global.DB, &param)
		}
		tenantId = param.TenantId
	case enum.ChangeStatus:
		var param form.RuleReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		if needSave {
			ruleDao.UpdateRuleState(global.DB, &param)
		}
		tenantId = param.TenantId
	case enum.DeleteRule:
		var param form.RuleReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		if needSave {
			ruleDao.DeleteRule(global.DB, &param)
		}
		tenantId = param.TenantId
	case enum.UnbindRule:
		var param form.UnBindRuleParam
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		if needSave {
			instanceDao.UnbindInstance(global.DB, &param)
		}
		tenantId = param.TenantId
	case enum.BindRule:
		var param form.InstanceBindRuleDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		if needSave {
			instanceDao.BindInstance(global.DB, &param)
		}
		tenantId = param.TenantId
	default:
		logger.Logger().Warnf("不支持的消息类型，消息类型：%v,消息%s", event, string(data))
	}
	if len(tenantId) > 0 {
		prometheusDao.GenerateUserPrometheusRule(tenantId)
	}
}
