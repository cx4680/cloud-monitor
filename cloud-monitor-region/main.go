package main

import (
	cp "code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/pipeline"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/pipeline"
	"context"
	"os"
)

func main() {

	c := context.Background()
	pl := cp.ActuatorPipeline{}
	if err := pl.First(&cp.ConfigActuatorStage{}).
		Then(&pipeline.TransactionActuatorStage{}).
		Then(&cp.SysComponentActuatorStage{}).
		Then(&pipeline.TaskActuatorStage{}).
		Then(&pipeline.MQActuatorStage{}).
		Then(&pipeline.K8sActuatorStage{}).
		Then(&pipeline.WebActuatorStage{}).
		Exec(&c); err != nil {
		os.Exit(1)
	}

}
