package form

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
)

type MqMsg struct {
	EventEum enum.EventEum
	Data     interface{}
}
