package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/handler_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"flag"
	"strconv"
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
	pipeline := (&ActuatorPipeline{}).First(func(c *context.Context) error {
		var cf = flag.String("config", "config.local.yml", "config.yml path")
		flag.Parse()
		return config.InitConfig(*cf)
	},
	).Then(func(c *context.Context) error {
		noticeChannelMap := make(map[string]string)
		if config.Cfg.Common.MsgIsOpen == config.MsgClose {
			global.NoticeChannelMap = noticeChannelMap
			return nil
		}
		msgChannelList := strings.Split(config.Cfg.Common.MsgChannel, ",")
		for _, v := range msgChannelList {
			switch v {
			case config.MsgChannelEmail:
				noticeChannelMap[v] = strconv.Itoa(handler_type.Email)
			case config.MsgChannelSms:
				noticeChannelMap[v] = strconv.Itoa(handler_type.Sms)
			}
		}
		global.NoticeChannelMap = noticeChannelMap
		return nil
	}).Then(func(c *context.Context) error {
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
