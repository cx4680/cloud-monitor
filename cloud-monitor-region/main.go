package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"

	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"flag"
	"log"
	"os"
)

func main() {
	//解析命令参数
	var cf = flag.String("config", "config.local.yml", "config.yml path")
	flag.Parse()

	//加载配置文件
	if err := config.InitConfig(*cf); err != nil {
		log.Printf("init config.yml error: %v\n", err)
		os.Exit(1)
	}

	if err := translate.InitTrans("zh"); err != nil {
		log.Printf("init trans failed, err:%v\n", err)
		os.Exit(2)
	}

	if err := sysComponent.InitSys(); err != nil {
		log.Printf("init sys component error: %v\n", err)
		os.Exit(3)
	}

	if config.GetCommonConfig().Env != "local" {
		if err := k8s.InitK8s(); err != nil {
			log.Printf("init k8s error: %v\n", err)
			os.Exit(5)
		}
	}
	//加载mq
	if err := initRocketMq(); err != nil {
		log.Printf("init rocketmq error, %v\n", err)
		os.Exit(5)
	}
	//加载定时任务
	go task.CronInstanceJob()

	logger.InitLogger(config.GetLogConfig())
	defer logger.Logger().Sync()

	//启动Web容器
	err := web.Start(config.GetServeConfig())
	if err != nil {
		logger.Logger().Infof("startup service failed, err:%v\n", err)
		os.Exit(6)
	}
}

func initRocketMq() error {
	rc := config.GetRocketmqConfig()
	if err := sysRocketMq.CreateTopics(rc.RuleTopic, rc.RecordTopic, rc.AlertContactTopic, rc.AlertContactGroup); err != nil {
		log.Printf("create topics error, %v\n", err)
		return err
	}

	if err := sysRocketMq.StartConsumersScribe([]sysRocketMq.Consumer{{
		Topic:   rc.AlertContactTopic,
		Handler: consumer.AlertContactHandler,
	}, {
		Topic:   rc.RuleTopic,
		Handler: consumer.AlarmRuleHandler,
	}}); err != nil {
		log.Printf("create rocketmq consumer error, %v\n", err)
		return err
	}
	return nil
}
