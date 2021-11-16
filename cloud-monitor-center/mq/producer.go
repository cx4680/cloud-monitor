package mq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var p rocketmq.Producer

func CreateMq() {
	cfg := config.GetRocketmqConfig()
	p, _ = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.NameServer})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return
	}
}

func SendMsg(topic string, eventEum enums.EventEum, module interface{}) error {

	var mqMsg = forms.MqMsg{
		EventEum: eventEum,
		Data:     module,
	}

	// msg对象转json ([]byte)
	jsonBytes, err := json.Marshal(mqMsg)
	if err != nil {
		fmt.Println(err)
	}

	msg := &primitive.Message{
		Topic: topic,
		Body:  jsonBytes,
	}
	res, err := p.SendSync(context.Background(), msg)

	if err != nil {
		return err
	} else {
		fmt.Printf("send message success: result=%s\n", res.String())
	}
	return nil
}
