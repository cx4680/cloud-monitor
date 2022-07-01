package sync

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/publisher"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/subscriber"
	"strings"
)

const defaultEmptyNameSrv = "${rocketMQEndpoint}"

func InitSync(regionRole string) error {
	if strings.EqualFold(regionRole, "integration") {
		logger.Logger().Info("current run mode: single region")
		publisher.GlobalPublisher = &publisher.NonePublisher{}
		return nil
	}
	logger.Logger().Info("current run mode: multiple regions")
	if isDefaultAddr() {
		logger.Logger().Info("rocketmq name server is default, don't init mq component")
		return nil
	}

	p, err := publisher.NewMQPublisher()
	if err != nil {
		logger.Logger().Errorf("init publisher error: %v\n", err)
		return err
	}
	publisher.GlobalPublisher = p

	sub := new(subscriber.Subscriber)
	if err := sub.Run(); err != nil {
		logger.Logger().Errorf("init subscriber error: %v\n", err)
		return err
	}
	return nil
}

func isDefaultAddr() bool {
	return config.Cfg.Rocketmq.NameServer == defaultEmptyNameSrv
}
