package actuator

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
)

type health struct {
	Status string `json:"status"`
}

func Health() string {
	return jsonutil.ToString(health{Status: "UP"})
}
