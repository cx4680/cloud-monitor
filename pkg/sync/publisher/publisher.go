package publisher

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"strings"
)

type PubMessage struct {
	Topic sys_rocketmq.Topic
	Data  interface{}
}

var GlobalPublisher Publisher

type Publisher interface {
	Pub(PubMessage) error
}

type NonePublisher struct{}

func (*NonePublisher) Pub(PubMessage) error {
	return nil
}

type MQPublisher struct {
	Producer rocketmq.Producer
}

func (mqp *MQPublisher) Pub(msg PubMessage) error {
	s, ok := msg.Data.(string)
	var body []byte
	if ok {
		body = []byte(s)
	} else {
		bs, err := json.Marshal(msg.Data)
		if err != nil {
			return err
		}
		body = bs
	}

	res, err := mqp.Producer.SendSync(context.Background(), &primitive.Message{
		Topic: string(msg.Topic),
		Body:  body,
	})
	if err != nil {
		return err
	}
	logger.Logger().Infof("publish MQ message, data=%v, res=%v", jsonutil.ToString(msg), res)
	return nil
}

func NewMQPublisher() (*MQPublisher, error) {
	p, err := initMQProducer()
	if err != nil {
		logger.Logger().Errorf("create rocketmq consumer error, %v\n", err)
		return nil, err
	}

	topics := []sys_rocketmq.Topic{
		sys_rocketmq.SmsMarginReminderTopic,
		sys_rocketmq.InstanceTopic,
		sys_rocketmq.ContactTopic,
		sys_rocketmq.ContactGroupTopic,
		sys_rocketmq.RecordTopic,
		sys_rocketmq.AlarmInfoTopic,
		sys_rocketmq.RuleTopic,
		sys_rocketmq.NotificationSyncTopic,
		sys_rocketmq.DeleteInstanceTopic,
		sys_rocketmq.MonitorProductTopic,
		sys_rocketmq.MonitorItemTopic,
	}

	err = sys_rocketmq.CreateTopics(topics...)
	if err != nil {
		logger.Logger().Errorf("create rocketmq topics error, %v\n", err)
		return nil, err
	}
	return &MQPublisher{Producer: p}, nil
}

func initMQProducer() (rocketmq.Producer, error) {
	cfg := config.Cfg.Rocketmq
	addrs := strings.Split(cfg.NameServer, ";")
	var err error
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(addrs)),
		producer.WithRetry(2),
	)
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return nil, err
	}
	if err := p.Start(); err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		return nil, err
	}
	return p, nil
}
