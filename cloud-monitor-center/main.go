package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
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
	var cf = flag.String("config", "config.local.yml", "config.yml path")
	flag.Parse()

	if err := config.InitConfig(*cf); err != nil {
		log.Printf("init config.yml error: %v\n", err)
		os.Exit(1)
	}

	if err := database.InitDb(config.GetDbConfig()); err != nil {
		log.Printf("init database error: %v\n", err)
		os.Exit(2)
	}

	if err := translate.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		os.Exit(3)
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
	if err := mq.CreateTopics(rc.RuleTopic, rc.RecordTopic, rc.AlertContactTopic, rc.AlertContactGroup); err != nil {
		log.Printf("create topics error, %v\n", err)
		return err
	}
	err := mq.InitProducer()
	if err != nil {
		log.Printf("create rocketmq producer error, %v\n", err)
		return err
	}
	//TODO 初始化消费者
	if err = mq.StartConsumersScribe([]mq.Consumer{{
		Topic:   "",
		Handler: nil,
	}}); err != nil {
		log.Printf("create rocketmq consumer error, %v\n", err)
		return err
	}
	return nil
}
