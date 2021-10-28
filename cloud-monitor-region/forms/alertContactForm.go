package forms

import "code.cestc.cn/ccos-ops/cloud-monitor/common/enums"

type MqMsg struct {
	EventEum enums.EventEum
	Data     interface{}
}
