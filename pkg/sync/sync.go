package sync

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/publisher"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/subscriber"
)

func InitSync(isSingle bool) error {
	if isSingle {
		publisher.GlobalPublisher = &publisher.NonePublisher{}
		return nil
	}

	p, err := publisher.NewMQPublisher()
	if err != nil {
		return err
	}
	publisher.GlobalPublisher = p

	sub := new(subscriber.Subscriber)
	if err := sub.Run(); err != nil {
		return err
	}
	return nil
}
