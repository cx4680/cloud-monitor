package forms

import "code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums"

type MqMsg struct {
	EventEum enums.EventEum
	Data     interface{}
}
