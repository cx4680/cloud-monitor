package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/web"
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
	var cf = flag.String("config.yml", "D:\\dev-go\\cloud-monitor\\cloud-monitor-center\\config.local.yml", "config.yml path")
	flag.Parse()

	//加载配置文件
	err := config.InitConfig(*cf)
	if err != nil {
		log.Printf("init config.yml error: %v\n", err)
		os.Exit(1)
	}

	database.InitDb(config.GetDbConfig())
	if err := translate.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	//创建mq
	mq.CreateMq()

	logger.InitLogger(config.GetLogConfig())
	defer logger.Logger().Sync()

	//启动Web容器
	if err := web.Start(config.GetServeConfig()); err != nil {
		logger.Logger().Infof("startup service failed, err:%v\n", err)
	}
}
