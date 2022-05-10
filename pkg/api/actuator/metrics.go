package actuator

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"runtime"
)

func Metrics() string {
	return jsonutil.ToString(refillMetricsMap())
}

func refillMetricsMap() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}
