package main

import (
	commonLoader "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysGuide"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/loader"
	"os"
)

// @title Swagger Example API
// @version 0.0.1
// @description  This is a sample server Petstore server.
// @BasePath /
func main() {

	l := commonLoader.SysSysGuideImpl{InitList: []commonLoader.SysLoader{
		&commonLoader.ConfigLoader{},
		&loader.TransactionLoader{},
		&commonLoader.SysComponentLoader{},
		&loader.RocketMQConsumerLoader{},
		&loader.TaskLoader{},
		&loader.WebServeLoader{},
	}}

	if err := l.StartServe(); err != nil {
		os.Exit(-1)
	}
}
