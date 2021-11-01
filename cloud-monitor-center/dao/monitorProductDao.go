package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/models"
	"gorm.io/gorm"
)

type MonitorProductDao struct {
	db *gorm.DB
}

func NewMonitorProductDao(db *gorm.DB) *MonitorProductDao {
	return &MonitorProductDao{db: db}
}

func (mpd *MonitorProductDao) SelectMonitorProductList() *[]models.MonitorProduct {
	var product = &[]models.MonitorProduct{}
	mpd.db.Where("status = ?", "1").Find(product)
	return product
}

func (mpd *MonitorProductDao) Create(product *models.MonitorProduct) {
	mpd.db.Create(product)
}

func (mpd *MonitorProductDao) GetById(id string) *models.MonitorProduct {
	var product models.MonitorProduct
	mpd.db.First(&product, id)
	return &product
}

func (mpd *MonitorProductDao) UpdateById(product *models.MonitorProduct) {
	mpd.db.Model(product).Updates(*product)
}

func (mpd *MonitorProductDao) DeleteById(id string) {
	var product models.MonitorProduct
	mpd.db.Delete(&product, id)
}
