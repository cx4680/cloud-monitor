package mq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func SendNotificationRecordMsg(msg []models.NotificationRecord) {
	if msg == nil || len(msg) <= 0 {
		return
	}
	doMqSendMsg("notification_sync", tools.ToString(msg))
}

func SendAlertRecordMsg(msg []*models.AlertRecord) {
	cfg := config.GetConfig()
	doMqSendMsg(cfg.Rocketmq.RecordTopic, tools.ToString(msg))
}

func doMqSendMsg(topic, msg string) {
	cfg := config.GetConfig()
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.Rocketmq.NameServer})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return
	}

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

	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}

}
