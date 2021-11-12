package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	//解析命令参数
	var cf = flag.String("config.yml", "config.local.yml", "config.yml path")
	flag.Parse()

	//加载配置文件
	err := config.InitConfig(*cf)
	if err != nil {
		log.Printf("init config.yml error: %v\n", err)
		os.Exit(1)
	}

	cfg := config.GetConfig()

	if err := translate.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	//初始化数据里连接
	database.InitDb(&cfg.DB)

	redisConfig := redis.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
	}
	redis.InitClient(redisConfig)
	k8s.InitK8s()
	//加载mq
	mq.SubScribe()
	//加载定时任务
	task.CronInstanceJob()

	logger.InitLogger(&config.GetConfig().Logger)
	defer logger.Logger().Sync()

	//启动Web容器
	err = web.Start(cfg)
	if err != nil {
		logger.Logger().Infof("startup service failed, err:%v\n", err)
	}
}
