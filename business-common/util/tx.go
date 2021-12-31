package util

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"fmt"
	"gorm.io/gorm"
)

func Tx(param interface{}, f func(xx *gorm.DB, param interface{}) error) (err error) {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			logger.Logger().Errorf("%v", err)
			tx.Rollback()
			err = fmt.Errorf("%v", err)
		}
	}()
	err = f(tx, param)
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
