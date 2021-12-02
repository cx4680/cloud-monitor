package main

import (
	commonLoader "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysGuide"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/loader"
	"os"
)

func main() {

	guide := commonLoader.SysSysGuideImpl{}
	guide.RegisterLoader(&commonLoader.ConfigLoader{})
	guide.RegisterLoader(&loader.TransactionLoader{})
	guide.RegisterLoader(&loader.K8sLoader{})
	guide.RegisterLoader(&commonLoader.SysComponentLoader{})
	guide.RegisterLoader(&loader.RocketMQConsumerLoader{})
	guide.RegisterLoader(&loader.TaskLoader{})
	guide.RegisterLoader(&loader.WebServeLoader{})

	if err := guide.StartServe(); err != nil {
		os.Exit(-1)
	}

}
