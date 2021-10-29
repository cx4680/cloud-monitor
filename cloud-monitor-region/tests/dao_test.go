package tests

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestMonitorProductDaoTest(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/hawkeye?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var list = []models.MonitorProduct{
		{

			Name:        "1",
			Status:      1,
			Description: "1",
			CreateUser:  "1",
			CreateTime:  "2021-01-01 00：00：00",
		},
		{
			Name:        "2",
			Status:      1,
			Description: "1",
			CreateUser:  "1",
			CreateTime:  "2021-01-01 00：00：00",
		},
	}

	db.Create(&list)
}
