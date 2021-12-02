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

	guide := commonLoader.SysSysGuideImpl{}
	guide.RegisterLoader(&commonLoader.ConfigLoader{})
	guide.RegisterLoader(&loader.TransactionLoader{})
	guide.RegisterLoader(&commonLoader.SysComponentLoader{})
	guide.RegisterLoader(&loader.RocketMQConsumerLoader{})
	guide.RegisterLoader(&loader.TaskLoader{})
	guide.RegisterLoader(&loader.WebServeLoader{})

	if err := guide.StartServe(); err != nil {
		os.Exit(1)
	}
}
