package form

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
)

type MonitorProductParam struct {
	BizIdList []string `form:"bizIdList"`
	Status    uint8    `form:"status"`
	EventEum  enum.EventEum
}
