package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"gorm.io/gorm"
)

type MonitorProductDao struct {
}

var MonitorProduct = new(MonitorProductDao)

func (mpd *MonitorProductDao) GetByAbbreviation(db *gorm.DB, abbreviation string) *models.MonitorProduct {
	if tools.IsBlank(abbreviation) {
		return nil
	}
	var product models.MonitorProduct
	db.Where(models.MonitorProduct{Abbreviation: abbreviation}).First(&product)
	return &product

}
func (mpd *MonitorProductDao) SelectMonitorProductList() *[]models.MonitorProduct {
	var product = &[]models.MonitorProduct{}
	global.DB.Where("status = ?", "1").Find(product)
	return product
}

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
