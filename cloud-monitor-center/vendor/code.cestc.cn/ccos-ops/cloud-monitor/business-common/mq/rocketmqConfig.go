package mq

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

type Consumer struct {
	Topic   string
	Handler func([]*primitive.MessageExt)
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

func InitProducer()  error {
	cfg := config.GetRocketmqConfig()
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.NameServer})),
		producer.WithRetry(2),
	)
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return  err
	}
	if err := p.Start(); err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return  err
	}
	return  nil
}

func StartConsumersScribe(consumers []Consumer) error {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(config.GetCommonConfig().RegionName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.GetRocketmqConfig().NameServer})),
	)
	for _, o := range consumers {
		err := c.Subscribe(o.Topic, consumer.MessageSelector{}, func(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			o.Handler(msg)
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

func SendMsg(topic, msg string) {
	mqmsg := &primitive.Message{
		Topic: topic,
		Body:  []byte(msg),
	}
	res, err := p.SendSync(context.Background(), mqmsg)

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	} else {
		fmt.Printf("send message success: result=%s\n", res.String())
	}
}

func Shutdown() {
	if err := p.Shutdown(); err != nil {
		log.Printf("shutdown rocketmq producer error, %v\n", err)
	}
	if err := c.Shutdown(); err != nil {
		log.Printf("shutdown rocketmq consumer error, %v\n", err)
	}
}
