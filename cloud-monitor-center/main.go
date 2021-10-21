package main

import (
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
	var cf = flag.String("config", "cloud-monitor-center/config.local.yml", "config path")
	flag.Parse()

	//加载配置文件
	err := config.InitConfig(*cf)
	if err != nil {
		log.Printf("init config error: %v\n", err)
		os.Exit(1)
	}

	cfg := config.GetConfig()

	if err := translate.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	//初始化数据里连接
	database.InitDb(&cfg.DB)
	logger.InitLogger(&config.GetConfig().Logger)
	defer logger.Logger().Sync()

	defer database.GetDb().Close()
	//启动Web容器
	if err := web.Start(cfg); err != nil {
		logger.Logger().Infof("startup service failed, err:%v\n", err)
	}
}
