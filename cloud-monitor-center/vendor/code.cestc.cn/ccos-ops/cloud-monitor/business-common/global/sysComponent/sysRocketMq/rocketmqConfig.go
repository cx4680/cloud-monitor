package sysRocketMq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/pkg/errors"
	"strings"
	"sync"
)

var p rocketmq.Producer

type Topic string
type Group string

const (
	SmsMarginReminderTopic Topic = "sms_margin_reminder" //短信余量提醒
	NotificationSyncTopic  Topic = "all_notification"    //通知记录

	RuleTopic     Topic = "alert_rule"                 //告警规则
	RecordTopic   Topic = "alert_record"               //告警历史记录
	InstanceTopic Topic = "alert_correlation_instance" //告警关联实例

	AlertContactTopic      Topic = "alert_contact"       //告警联系人
	AlertContactGroupTopic Topic = "alert_contact_group" //告警联系组

	DeleteInstanceTopic Topic = "delete_instance"
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
	cfg := config.GetRocketmqConfig()
	addrs := strings.Split(cfg.NameServer, ";")
	brokerAddr, err2 := getBrokerAddr(cfg)
	if err2 != nil {
		return err2
	}

	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(addrs)))
	if err != nil {
		logger.Logger().Infof("create topic error %+v\n", err)
		return err
	}
	for _, topic := range topics {
		err = testAdmin.CreateTopic(
			context.Background(),
			admin.WithTopicCreate(string(topic)),
			admin.WithBrokerAddrCreate(brokerAddr),
		)
		if err != nil {
			logger.Logger().Infof("create topic error %+v", err)
			return err
		}
	}
	return nil
}

func getBrokerAddr(cfg config.Rocketmq) (string, error) {
	c := &Client{
		responseTable:    sync.Map{},
		connectionTable:  sync.Map{},
		option:           TcpOption{},
		processors:       nil,
		connectionLocker: sync.Mutex{},
		interceptor:      nil,
	}
	broker, e := c.GetBrokerDataList(strings.Split(cfg.NameServer, ";")[0])
	if e != nil {
		return "", e
	}
	brokerAddr := ""
	for _, b := range broker.BrokerDataList {
		for _, v := range b.BrokerAddresses {
			brokerAddr = v
			break
		}
	}
	return brokerAddr, nil
}

func InitProducer() error {
	cfg := config.GetRocketmqConfig()
	addrs := strings.Split(cfg.NameServer, ";")
	var err error
	p, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(addrs)),
		producer.WithRetry(2),
	)
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return err
	}
	if err := p.Start(); err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return err
	}
	return nil
}

func StartConsumersScribe(group Group, consumers []*Consumer) error {
	addrs := strings.Split(config.GetRocketmqConfig().NameServer, ";")
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(string(group)),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(addrs)),
	)
	for _, o := range consumers {
		m := *o
		err := c.Subscribe(string(m.Topic), consumer.MessageSelector{}, func(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
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

func SendRocketMqMsg(msg RocketMqMsg) error {
	return SendMsg(msg.Topic, msg.Content)
}

func SendMsg(topic Topic, msg string) error {
	if topic == "" {
		return errors.New("topic can't be null")
	}
	if msg == "" {
		return errors.New("rocketmq send msg can't be null")
	}
	res, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: string(topic),
		Body:  []byte(msg),
	})

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
		return err
	}
	fmt.Printf("send message success: result=%s\n", res.String())
	return nil
}
