package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/models"
	"gorm.io/gorm"
)

type MonitorItemDao struct {
	db *gorm.DB
}

func NewMonitorItemDao() *MonitorItemDao {
	return &MonitorItemDao{db: database.GetDb()}
}

func (mpd *MonitorItemDao) GetLabelsByName(name string) string {
	var model = models.MonitorItem{}
	mpd.db.Debug().Where("metric_name = ?", name).First(&model)
	return model.Labels
}
