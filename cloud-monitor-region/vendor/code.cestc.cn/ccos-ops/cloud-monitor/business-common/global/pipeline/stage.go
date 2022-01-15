package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"flag"
	"strings"
)

type SysLoader interface {
	AddStage(Actuator) SysLoader
	Start() (*context.Context, error)
}

type MainLoader struct {
	Pipeline Pipeline
}

func NewMainLoader() *MainLoader {
	pipeline := (&ActuatorPipeline{}).First(func(*context.Context) error {
		var cf = flag.String("config", "config.local.yml", "config.yml path")
		flag.Parse()
		return config.InitConfig(*cf)
	},
	).Then(func(*context.Context) error {
		if config.Cfg.Common.MsgIsOpen == config.MsgClose {
			return nil
		}
		msgChannelList := strings.Split(config.Cfg.Common.MsgChannel, ",")
		for _, v := range msgChannelList {
			switch v {
			case config.MsgChannelEmail:
				global.NoticeChannelList = append(global.NoticeChannelList, form.NoticeChannel{Name: "邮箱", Code: v, Data: 1})
			case config.MsgChannelSms:
				global.NoticeChannelList = append(global.NoticeChannelList, form.NoticeChannel{Name: "短信", Code: v, Data: 2})
			}
		}
		return nil
	}).Then(func(*context.Context) error {
		return sys_component.InitSys()
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
