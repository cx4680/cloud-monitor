package main

import (
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/pipeline"
	"fmt"
	"os"
)

func main() {

	loader := cp.NewMainLoader().
		AddStage(&pipeline.IamActuatorStage{}).
		AddStage(&pipeline.DBInitActuatorStage{}).
		AddStage(&pipeline.TransactionActuatorStage{}).
		AddStage(&pipeline.TaskActuatorStage{}).
		AddStage(&pipeline.K8sActuatorStage{}).
		AddStage(&pipeline.MQActuatorStage{}).
		AddStage(&pipeline.WebActuatorStage{})

	_, err := loader.Start()
	if err != nil {
		fmt.Printf("exit error:%v", err)
		os.Exit(1)
	}
}
