package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"flag"
)

type SysLoader interface {
	AddStage(Actuator) SysLoader
	Start() (*context.Context, error)
}

type MainLoader struct {
	Pipeline Pipeline
}

func NewMainLoader() *MainLoader {
	pipeline := (&ActuatorPipeline{}).First(func(c *context.Context) error {
		var cf = flag.String("config", "config.local.yml", "config.yml path")
		flag.Parse()
		return config.InitConfig(*cf)
	},
	).Then(func(c *context.Context) error {
		return sysComponent.InitSys()
	})
	return &MainLoader{Pipeline: pipeline}
}

func (l *MainLoader) AddStage(actuator Actuator) SysLoader {
	l.Pipeline = l.Pipeline.Then(actuator)
	return l
}

func (l *MainLoader) Start() (*context.Context, error) {
	c := context.Background()
	return &c, l.Pipeline.Exec(&c)
}
