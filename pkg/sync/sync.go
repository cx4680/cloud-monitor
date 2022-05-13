package sync

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/publisher"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/subscriber"
	"strings"
)

func InitSync(isSingle string) error {
	if len(isSingle) == 0 || strings.EqualFold(isSingle, "true") {
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
