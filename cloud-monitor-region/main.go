package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	//解析命令参数
	var cf = flag.String("config", "config.local.yml", "config path")
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
	defer logger.Logger.Sync()
	logger.Logger.Infof("xxxxx123", 123)

	zap.L().Error("ssss123")

	defer database.GetDb().Close()
	//启动Web容器
	if err := web.Start(cfg); err != nil {
		logger.Logger.Info("startup service failed, err:%v\n", err)
	}
}
