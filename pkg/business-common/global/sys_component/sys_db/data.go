package sys_db

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
)

var Update = false

type migrateLog struct {
}

func (l *migrateLog) Printf(format string, v ...interface{}) {
	logger.Logger().Infof(format, v...)
}

func (l *migrateLog) Verbose() bool {
	return true
}
func InitData(dbConfig config.DB, database, path string) error {
	if !global.DB.Migrator().HasTable(&model.AlarmHandler{}) && global.DB.Migrator().HasColumn(&model.AlarmRule{}, "notify_channel") {
		Update = true
	}
	//升级到 5_312 需要重新生成 prometheus rule
	if !Update && !global.DB.Migrator().HasTable(&model.AlarmItem{}) && global.DB.Migrator().HasTable(&model.AlarmHandler{}) {
		Update = true
	}
	var err error
	db, err := global.DB.DB()
	if err != nil {
		return err
	}
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(path, database, driver)
	if err != nil {
		return err
	}
	m.Log = new(migrateLog)
	m.LockTimeout = 3 * time.Minute
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Logger().Error("An error occurred while syncing the database.. ", err)
		return err
	}
	return nil
}
