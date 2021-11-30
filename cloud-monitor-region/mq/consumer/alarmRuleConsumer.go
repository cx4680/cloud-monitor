package consumer

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonForm "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"

	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func AlarmRuleHandler(msgList []*primitive.MessageExt) {
	for i := range msgList {
		fmt.Printf("subscribe callback: %v \n", msgList[i])
		if err := json.Unmarshal(msgList[i].Body, &MqMsg); err != nil {
			continue
		}
		data, err := json.Marshal(MqMsg.Data)
		if err != nil {
			continue
		}
		err = tx(MqMsg, data)
		if err != nil {
			fmt.Printf("sync error: %v ,%v\n", msgList[i], err)
		}
	}
}

func tx(MqMsg forms.MqMsg, data []byte) error {
	tx := global.DB
	var err error
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Errorf("%v", err)
			tx.Rollback()
			err = fmt.Errorf("%v", err)
		}
	}()
	handleMsg(MqMsg, data)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}

func handleMsg(MqMsg forms.MqMsg, data []byte) {
	ruleDao := commonDao.AlarmRule
	instanceDao := commonDao.Instance
	prometheusDao := service.PrometheusRule
	var tenantId string
	switch MqMsg.EventEum {
	case enums.CreateRule:
		var param commonForm.AlarmRuleAddReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.SaveRule(global.DB, &param)
		tenantId = param.TenantId
	case enums.UpdateRule:
		var param commonForm.AlarmRuleAddReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.UpdateRule(global.DB, &param)
		tenantId = param.TenantId
	case enums.EnableRule:
		var param commonForm.RuleReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.UpdateRuleState(global.DB, &param)
	case enums.DisableRule:
		var param commonForm.RuleReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.UpdateRuleState(global.DB, &param)
	case enums.DeleteRule:
		var param commonForm.RuleReqDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		ruleDao.DeleteRule(global.DB, &param)
	case enums.UnbindRule:
		var param commonForm.UnBindRuleParam
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		instanceDao.UnbindInstance(global.DB, &param)
	case enums.BindRule:
		var param commonForm.InstanceBindRuleDTO
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		instanceDao.BindInstance(global.DB, &param)
	default:
		logger.Logger().Warnf("不支持的消息类型，消息类型：%v,消息%s", MqMsg.EventEum, string(data))
	}
	if len(tenantId) > 0 {
		prometheusDao.GenerateUserPrometheusRule(tenantId)
	}
}
