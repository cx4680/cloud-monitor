package actuator

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"runtime"
)

func Metrics() string {
	return tools.ToString(refillMetricsMap())
}

func refillMetricsMap() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}
