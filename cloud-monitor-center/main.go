package main

import (
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/pipeline"
	"context"
	"os"
)

// @title Swagger Example API
// @version 0.0.1
// @description  This is a sample server Petstore server.
// @BasePath /
func main() {
	c := context.Background()
	pl := cp.ActuatorPipeline{}
	if err := pl.First(&cp.ConfigActuatorStage{}).
		Then(&pipeline.TransactionActuatorStage{}).
		Then(&cp.SysComponentActuatorStage{}).
		Then(&pipeline.TaskActuatorStage{}).
		Then(&pipeline.MQActuatorStage{}).
		Then(&pipeline.WebActuatorStage{}).
		Exec(&c); err != nil {
		os.Exit(1)
	}
}
