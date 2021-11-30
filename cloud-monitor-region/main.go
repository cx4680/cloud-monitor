package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
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
	if err := initRocketMqConsumers(); err != nil {
		log.Printf("init rocketmq error, %v\n", err)
		os.Exit(5)
	}
	//加载定时任务
	if err := initTask(); err != nil {
		log.Printf("init task error, %v\n", err)
		os.Exit(6)
	}

	logger.InitLogger(config.GetLogConfig())
	defer logger.Logger().Sync()

	//启动Web容器
	err := web.Start(config.GetServeConfig())
	if err != nil {
		logger.Logger().Infof("startup service failed, err:%v\n", err)
		os.Exit(6)
	}
}

func initRocketMqConsumers() error {
	if err := sysRocketMq.StartConsumersScribe(sysRocketMq.Group(config.GetCommonConfig().RegionName), []*sysRocketMq.Consumer{{
		Topic:   sysRocketMq.AlertContactTopic,
		Handler: consumer.AlertContactHandler,
	}, {
		Topic:   sysRocketMq.RuleTopic,
		Handler: consumer.AlarmRuleHandler,
	}}); err != nil {
		log.Printf("create rocketmq consumer error, %v\n", err)
		return err
	}
	return nil
}

func initTask() error {
	bt := commonTask.NewBusinessTaskImpl()
	if err := bt.Add(commonTask.BusinessTaskDTO{
		Cron: "0 0 0/1 * * ?",
		Name: "instanceJob",
		Task: task.InstanceJob,
	}); err != nil {
		return err
	}

	if err := bt.Add(commonTask.BusinessTaskDTO{
		Cron: "0 0 0/1 * * ?",
		Task: task.SlbInstanceJob,
	}); err != nil {
		return err
	}

	bt.Start()
	return nil
}
