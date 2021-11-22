package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
)

type MonitorProductDao struct {
}

var MonitorProduct = new(MonitorProductDao)

func (mpd *MonitorProductDao) SelectMonitorProductList() *[]models.MonitorProduct {
	var product = &[]models.MonitorProduct{}
	database.GetDb().Where("status = ?", "1").Find(product)
	return product
}

func (mpd *MonitorProductDao) Create(product *models.MonitorProduct) {
	database.GetDb().Create(product)
}

func (mpd *MonitorProductDao) GetById(id string) *models.MonitorProduct {
	var product models.MonitorProduct
	database.GetDb().First(&product, id)
	return &product
}

func (mpd *MonitorProductDao) UpdateById(product *models.MonitorProduct) {
	database.GetDb().Model(product).Updates(*product)
}

func (mpd *MonitorProductDao) DeleteById(id string) {
	var product models.MonitorProduct
	database.GetDb().Delete(&product, id)
}
