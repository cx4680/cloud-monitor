package dbUtils

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"fmt"
	"gorm.io/gorm"
)

func Tx(param interface{}, dao interface{}, f func(xx *gorm.DB, param interface{}, dao interface{}) error) (err error) {
	tx := database.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Errorf("%v", err)
			tx.Rollback()
			err = fmt.Errorf("%v", err)
		}
	}()
	err = f(tx, param, dao)
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return
	}
	return err
}
