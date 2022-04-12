package sys_rocketmq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"strings"
	"sync"
)

type Topic string
type Group string

const (
	SmsMarginReminderTopic Topic = "sms_margin_reminder" //短信余量提醒
	NotificationSyncTopic  Topic = "all_notification"    //通知记录

	RuleTopic  Topic = "alert_rule" //告警规则
	AlarmTopic       = "alarm"      //告警数据同步：告警历史+告警详情

	InstanceTopic Topic = "alert_correlation_instance" //告警关联实例

	ContactTopic      Topic = "contact"       //告警联系人
	ContactGroupTopic Topic = "contact_group" //告警联系组

	DeleteInstanceTopic Topic = "delete_instance"

	MonitorProductTopic Topic = "monitor_product" //监控产品
	MonitorItemTopic    Topic = "monitor_item"    //监控项
)

type Consumer struct {
	Topic   Topic
	Handler func([]*primitive.MessageExt)
}

type RocketMqMsg struct {
	Topic   Topic
	Content string
}

func CreateTopics(topics ...Topic) error {
	cfg := config.Cfg.Rocketmq
	addrs := strings.Split(cfg.NameServer, ";")
	brokerAddrs, err2 := getBrokerAddrs(cfg)
	if err2 != nil {
		return err2
	}

	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(addrs)))
	if err != nil {
		logger.Logger().Infof("create topic error %+v\n", err)
		return err
	}
	for _, topic := range topics {
		for _, addr := range brokerAddrs {
			err = testAdmin.CreateTopic(
				context.Background(),
				admin.WithTopicCreate(string(topic)),
				admin.WithBrokerAddrCreate(addr),
			)
			if err != nil {
				logger.Logger().Infof("create topic error %+v", err)
				return err
			}
		}
	}

	return nil
}

func getBrokerAddrs(cfg config.Rocketmq) ([]string, error) {
	c := &Client{
		responseTable:    sync.Map{},
		connectionTable:  sync.Map{},
		option:           TcpOption{},
		processors:       nil,
		connectionLocker: sync.Mutex{},
		interceptor:      nil,
	}
	broker, e := c.GetBrokerDataList(cfg.NameServer)
	if e != nil {
		logger.Logger().Info("get broker address error", e)
		return nil, e
	}
	var brokerAddrs []string
	for _, b := range broker.BrokerDataList {
		for _, v := range b.BrokerAddresses {
			brokerAddrs = append(brokerAddrs, v)
		}
	}
	return brokerAddrs, nil
}

func StartConsumersScribe(group Group, consumers []*Consumer) error {
	addresses := strings.Split(config.Cfg.Rocketmq.NameServer, ";")
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(string(group)),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(addresses)),
	)
	for _, o := range consumers {
		m := *o
		err := c.Subscribe(string(m.Topic), consumer.MessageSelector{
			Type:       consumer.SQL92,
			Expression: "TAGS <> '" + config.Cfg.Common.RegionName + "'",
		}, func(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			m.Handler(msg)
			return consumer.ConsumeSuccess, nil
		})
		if err != nil {
			return err
		}
	}
	if err := c.Start(); err != nil {
		return err
	}
	return nil
}
