package sysRocketMq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
)

var p rocketmq.Producer
var c rocketmq.PushConsumer

type Topic string

const (
	SmsMarginReminder Topic = "sms_margin_reminder" //短信余量提醒

	NotificationSync Topic = "notification_sync" //通知记录

)

type Consumer struct {
	Topic   string
	Handler func([]*primitive.MessageExt)
}

type RocketMqMsg struct {
	Topic   Topic
	Content string
}

func CreateTopics(topics ...string) error {
	cfg := config.GetRocketmqConfig()
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{cfg.NameServer})))
	if err != nil {
		fmt.Sprintf("create topic error %+v", err)
		return err
	}
	for _, topic := range topics {
		err = testAdmin.CreateTopic(
			context.Background(),
			admin.WithTopicCreate(topic),
			admin.WithBrokerAddrCreate(cfg.BrokerAddr),
		)
		if err != nil {
			fmt.Sprintf("create topic error %+v", err)
			return err
		}
	}
	return nil
}

func InitProducer() error {
	cfg := config.GetRocketmqConfig()
	var err error
	p, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.NameServer})),
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

func StartConsumersScribe(consumers []*Consumer) error {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(config.GetCommonConfig().RegionName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.GetRocketmqConfig().NameServer})),
	)
	for _, o := range consumers {
		m := *o
		err := c.Subscribe(m.Topic, consumer.MessageSelector{}, func(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
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
	//TODO
	return SendMsg(string(msg.Topic), msg.Content)
}

func SendMsg(topic, msg string) error {
	mqmsg := &primitive.Message{
		Topic: topic,
		Body:  []byte(msg),
	}
	res, err := p.SendSync(context.Background(), mqmsg)

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
		return err
	} else {
		fmt.Printf("send message success: result=%s\n", res.String())
	}
	return nil
}

func Shutdown() {
	if err := p.Shutdown(); err != nil {
		log.Printf("shutdown rocketmq producer error, %v\n", err)
	}
	if err := c.Shutdown(); err != nil {
		log.Printf("shutdown rocketmq consumer error, %v\n", err)
	}
}
