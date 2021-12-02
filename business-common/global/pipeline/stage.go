package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"flag"
)

type ConfigActuatorStage struct {
}

func (ca *ConfigActuatorStage) Exec(c *context.Context) error {
	var cf = flag.String("config", "config.local.yml", "config.yml path")
	flag.Parse()
	return config.InitConfig(*cf)
}

type SysComponentActuatorStage struct {
}

func (sa *SysComponentActuatorStage) Exec(c *context.Context) error {
	return sysComponent.InitSys()
}
