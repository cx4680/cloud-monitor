package database

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb(dbConfig *config.DB) {
	d, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbConfig.Url, // DSN data source name
		DefaultStringSize:         256,          // string 类型字段的默认长度
		DisableDatetimePrecision:  true,         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,        // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	/*// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	d.SetMaxIdleConns(dbConfig.MaxIdleConnes)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConnes)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(dbConfig.MaxLifeTime)
	d.LogMode(true)
	d.SingularTable(true)
	db = d*/
	db = d
}

func GetDb() *gorm.DB {
	return db
}
