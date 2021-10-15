package database

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/config"
	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/mysql"
)

var db *gorm.DB

func InitDb(dbConfig *config.DB) {
	d, err := gorm.Open(dbConfig.Dialect, dbConfig.Url)
	if err != nil {
		panic(err)
	}
	sqlDB := d.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConnes)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConnes)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(dbConfig.MaxLifeTime)
	d.SingularTable(true)
	db = d
}

func GetDb() *gorm.DB {
	return db
}
