package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
)

type MonitorProductDao struct {
}

var MonitorProduct = new(MonitorProductDao)

func (mpd *MonitorProductDao) Create(product *models.MonitorProduct) {
	global.DB.Create(product)
}

func (mpd *MonitorProductDao) GetById(id string) *models.MonitorProduct {
	var product models.MonitorProduct
	global.DB.First(&product, id)
	return &product
}

func (mpd *MonitorProductDao) UpdateById(product *models.MonitorProduct) {
	global.DB.Model(product).Updates(*product)
}

func (mpd *MonitorProductDao) DeleteById(id string) {
	var product models.MonitorProduct
	global.DB.Delete(&product, id)
}
