package main

import (
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/pipeline"
	"os"
)

func main() {

	loader := cp.NewMainLoader().AddStage(&pipeline.TransactionActuatorStage{}).
		AddStage(&pipeline.TaskActuatorStage{}).
		AddStage(&pipeline.MQActuatorStage{}).
		AddStage(&pipeline.K8sActuatorStage{}).
		AddStage(&pipeline.WebActuatorStage{})

	_, err := loader.Start()
	if err != nil {
		os.Exit(1)
	}
}
