package consumer

import (
	dao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	forms2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"golang.org/x/net/context"
)

var ruleDao *dao.AlarmRuleDao
var instanceDao *dao.InstanceDao
var prometheusDao *dao2.PrometheusRuleDao
var ruleConsumer *rocketmq.PushConsumer

func AlarmRuleConsumer() {
	cfg := config.GetRocketmqConfig()
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(cfg.RuleTopic),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.NameServer})),
	)

	ruleDao = dao.NewAlarmRuleDao(database.GetDb())
	instanceDao = dao.NewInstanceDao(database.GetDb())
	prometheusDao = dao2.NewPrometheusRuleDao(database.GetDb())
	var MqMsg forms.MqMsg

	err := c.Subscribe(cfg.RuleTopic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
			json.Unmarshal(msgs[i].Body, &MqMsg)
			data, _ := json.Marshal(MqMsg.Data)
			var tenantId string
			tx(MqMsg, data, tenantId)
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ruleConsumer = &c
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
	switch MqMsg.EventEum {
	case enums.CreateRule:
		var param forms2.AlarmRuleAddReqDTO
		json.Unmarshal(data, &param)
		ruleDao.SaveRule(ruleDao.DB, &param)
		tenantId = param.TenantId
	case enums.UpdateRule:
		var param forms2.AlarmRuleAddReqDTO
		json.Unmarshal(data, &param)
		ruleDao.UpdateRule(ruleDao.DB, &param)
		tenantId = param.TenantId
	case enums.EnableRule:
		var param forms2.RuleReqDTO
		json.Unmarshal(data, &param)
		ruleDao.UpdateRuleState(ruleDao.DB, &param)
	case enums.DisableRule:
		var param forms2.RuleReqDTO
		json.Unmarshal(data, &param)
		ruleDao.UpdateRuleState(ruleDao.DB, &param)
	case enums.DeleteRule:
		var param forms2.RuleReqDTO
		json.Unmarshal(data, &param)
		ruleDao.DeleteRule(ruleDao.DB, &param)
	case enums.UnbindRule:
		var param forms2.UnBindRuleParam
		json.Unmarshal(data, &param)
		instanceDao.UnbindInstance(ruleDao.DB, &param)
	case enums.BindRule:
		var param forms2.InstanceBindRuleDTO
		json.Unmarshal(data, &param)
		instanceDao.BindInstance(ruleDao.DB, &param)
	default:
		logger.Logger().Warnf("不支持的消息类型，消息类型：%v,消息%s", MqMsg.EventEum, string(data))
	}
	if len(tenantId) > 0 {
		prometheusDao.GenerateUserPrometheusRule("", "", tenantId)
	}
}

func ShowDown() {
	if ruleConsumer != nil {
		err := (*ruleConsumer).Shutdown()
		if err != nil {
			fmt.Printf("shutdown Consumer error: %s", err.Error())
		}
	}
}
