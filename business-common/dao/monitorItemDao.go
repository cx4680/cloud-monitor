package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
)

type MonitorItemDao struct {
}

var MonitorItem = new(MonitorItemDao)

func (mpd *MonitorItemDao) GetLabelsByName(name string) string {
	var model = models.MonitorItem{}
	global.DB.Debug().Where("metric_name = ?", name).First(&model)
	return model.Labels
}
