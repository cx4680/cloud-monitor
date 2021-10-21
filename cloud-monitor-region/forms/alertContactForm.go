package forms

import "code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/enums"

type MqMsg struct {
	EventEum enums.EventEum
	Data     interface{}
}
