package run_time

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
)

func SafeRun(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger().Error("run time error, ", err)
		}
	}()
	fn()
}
