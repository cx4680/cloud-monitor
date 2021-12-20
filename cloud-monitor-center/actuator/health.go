package actuator

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
)

type health struct {
	Status string `json:"status"`
}

func Health() string {
	return tools.ToString(health{Status: "UP"})
}
