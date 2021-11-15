package mq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func SubScribe() {
	consumer.AlertContactConsumer()
	consumer.AlarmRuleConsumer()
}

func Init() {
	cfg := config.GetRocketmqConfig()
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{cfg.NameServer})))
	if err != nil {
		fmt.Sprintf("create topic error %+v", err)
	}
	err = testAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(cfg.RuleTopic),
		admin.WithBrokerAddrCreate(cfg.BrokerAddr),
	)
	if err != nil {
		fmt.Sprintf("create topic error %+v", err)
	}
	err = testAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(cfg.AlertContactTopic),
		admin.WithBrokerAddrCreate(cfg.BrokerAddr),
	)
	if err != nil {
		fmt.Sprintf("create topic error %+v", err)
	}
}
