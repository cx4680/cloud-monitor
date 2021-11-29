package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"flag"
	"fmt"
	"log"
	"os"
)

// @title Swagger Example API
// @version 0.0.1
// @description  This is a sample server Petstore server.
// @BasePath /
func main() {
	//解析命令参数
	var cf = flag.String("config", "cloud-monitor-center/config.local.yml", "config.yml path")
	flag.Parse()

	if err := config.InitConfig(*cf); err != nil {
		log.Printf("init config.yml error: %v\n", err)
		os.Exit(1)
	}

	if err := translate.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		os.Exit(3)
	}

	if err := sysComponent.InitSys(); err != nil {
		fmt.Printf("init sys components failed, err:%v\n", err)
		os.Exit(2)
	}

	if err := initRocketMq(); err != nil {
		log.Printf("create rocketmq consumer error, %v\n", err)
		os.Exit(5)
	}

	logger.InitLogger(config.GetLogConfig())
	defer logger.Logger().Sync()

	//启动Web容器
	if err := web.Start(config.GetServeConfig()); err != nil {
		logger.Logger().Infof("startup service failed, err:%v\n", err)
		os.Exit(4)
	}
}

func initRocketMq() error {
	rc := config.GetRocketmqConfig()
	if err := sysRocketMq.CreateTopics(rc.RuleTopic, rc.RecordTopic, rc.AlertContactTopic, rc.InstanceTopic); err != nil {
		log.Printf("create topics error, %v\n", err)
		return err
	}
	//TODO 初始化消费者
	if err := sysRocketMq.StartConsumersScribe([]*sysRocketMq.Consumer{{
		Topic:   rc.InstanceTopic,
		Handler: consumer.InstanceHandler,
	}}); err != nil {
		log.Printf("create rocketmq consumer error, %v\n", err)
		return err
	}
	return nil
}

func initTask() error {
	//TODO
	return nil
}
