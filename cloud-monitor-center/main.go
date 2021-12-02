package main

import (
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/pipeline"
	"os"
)

// @title Swagger Example API
// @version 0.0.1
// @description  This is a sample server Petstore server.
// @BasePath /
func main() {

	loader := cp.NewMainLoader()
	loader.AddStage(&pipeline.TransactionActuatorStage{})
	loader.AddStage(&pipeline.TaskActuatorStage{})
	loader.AddStage(&pipeline.MQActuatorStage{})
	loader.AddStage(&pipeline.WebActuatorStage{})

	_, err := loader.Start()
	if err != nil {
		os.Exit(1)
	}

}
