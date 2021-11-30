package loader

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"flag"
)

type SysLoader interface {
	//StartServe 启动服务
	StartServe(loader SysLoader) error
	// ParseParameters 解析参数
	ParseParameters() string
	// InitConfig 加载配置文件
	InitConfig(file string) error
	// InitTrans 初始化翻译组件
	InitTrans() error
	// InitSysComponents 初始化系统组件
	InitSysComponents() error
	// InitRocketMqConsumers 初始化mq consumers
	InitRocketMqConsumers() error
	// InitTask 初始化定时任务
	InitTask() error
	// InitWebServe 初始化web容器
	InitWebServe() error
}

type AbstractSysLoader struct {
}

func (s *AbstractSysLoader) ParseParameters() string {
	var cf = flag.String("config", "config.local.yml", "config.yml path")
	flag.Parse()
	return *cf
}

func (s *AbstractSysLoader) StartServe(loader SysLoader) error {
	if err := s.InitConfig(s.ParseParameters()); err != nil {
		return err
	}

	if err := loader.InitTrans(); err != nil {
		return err
	}

	if err := s.InitSysComponents(); err != nil {
		return err
	}

	if err := loader.InitRocketMqConsumers(); err != nil {
		return err
	}
	if err := loader.InitTask(); err != nil {
		return err
	}

	if err := loader.InitWebServe(); err != nil {
		return err
	}

	return nil
}

func (s *AbstractSysLoader) InitConfig(file string) error {
	return config.InitConfig(file)
}

func (s *AbstractSysLoader) InitSysComponents() error {
	return sysComponent.InitSys()
}
