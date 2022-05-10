package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"gorm.io/gorm"
)

type AlarmInfoDao struct {
}

var AlarmInfo = new(AlarmInfoDao)

func (d *AlarmInfoDao) InsertBatch(db *gorm.DB, list []model.AlarmInfo) {
	db.Create(&list)
}
