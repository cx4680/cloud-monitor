package main

import (
	commonLoader "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysGuide"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/loader"
	"os"
)

func main() {

	l := commonLoader.SysSysGuideImpl{LoaderList: []commonLoader.SysLoader{
		&commonLoader.ConfigLoader{},
		&loader.TransactionLoader{},
		&loader.K8sLoader{},
		&commonLoader.SysComponentLoader{},
		&loader.RocketMQConsumerLoader{},
		&loader.TaskLoader{},
		&loader.WebServeLoader{},
	}}

	if err := l.StartServe(); err != nil {
		os.Exit(-1)
	}

}
