package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/pkg/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/pkg/database"
	"log"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {

	err := config.InitConfig("../config.dev.yml")
	if err != nil {
		log.Printf("init config error: %v\n", err)
		os.Exit(1)
	}

	database.InitDb(&config.GetConfig().DB)

	monitorProductDao := NewMonitorProductDao(database.GetDb())
	product := monitorProductDao.GetById("1")
	log.Println(product.Name)
}
