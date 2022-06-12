package sync

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/publisher"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/subscriber"
	"strings"
)

func InitSync(regionRole string) error {
	if strings.EqualFold(regionRole, "integration") {
		logger.Logger().Info("current run mode: single region")
		publisher.GlobalPublisher = &publisher.NonePublisher{}
		return nil
	}
	logger.Logger().Info("current run mode: multiple regions")

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
