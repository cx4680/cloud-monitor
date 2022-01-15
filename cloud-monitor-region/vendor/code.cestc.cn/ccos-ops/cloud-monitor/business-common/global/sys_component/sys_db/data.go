package sys_db

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"os"
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
	if !global.DB.Migrator().HasTable(&commonModels.AlarmHandler{}) && global.DB.Migrator().HasColumn(&commonModels.AlarmRule{}, "notify_channel") {
		Update = true
	}
	var err error
	pwd := os.Getenv("DB_PWD")
	url := dbConfig.Username + ":" + pwd + "@" + dbConfig.Url + "&multiStatements=true"
	db, err := sql.Open("mysql", url)
	defer db.Close()
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
