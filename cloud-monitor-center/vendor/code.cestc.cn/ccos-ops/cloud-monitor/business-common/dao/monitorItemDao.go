package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/models"
)

type MonitorItemDao struct {
}

var MonitorItem = new(MonitorItemDao)

func (mpd *MonitorItemDao) GetLabelsByName(name string) string {
	var model = models.MonitorItem{}
	database.GetDb().Debug().Where("metric_name = ?", name).First(&model)
	return model.Labels
}
