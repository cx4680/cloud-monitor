package sysGuide

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"flag"
)

type SysLoader interface {
	Load() error
}

type SysGuide interface {
	StartServe() error
}

type SysSysGuideImpl struct {
	loaderList []SysLoader
}

func (s *SysSysGuideImpl) RegisterLoader(l SysLoader) {
	s.loaderList = append(s.loaderList, l)
}

func (s *SysSysGuideImpl) StartServe() error {
	for _, init := range s.loaderList {
		if err := init.Load(); err != nil {
			return err
		}
	}
	return nil
}

type ConfigLoader struct{}

func (p *ConfigLoader) Load() error {
	var cf = flag.String("config", "config.local.yml", "config.yml path")
	flag.Parse()
	return config.InitConfig(*cf)
}

type SysComponentLoader struct{}

func (p *SysComponentLoader) Load() error {
	return sysComponent.InitSys()
}
