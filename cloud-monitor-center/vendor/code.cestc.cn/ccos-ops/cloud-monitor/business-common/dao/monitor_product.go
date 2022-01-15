package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"gorm.io/gorm"
)

type MonitorProductDao struct {
}

var MonitorProduct = new(MonitorProductDao)

func (mpd *MonitorProductDao) GetByAbbreviation(db *gorm.DB, abbreviation string) *model.MonitorProduct {
	if strutil.IsBlank(abbreviation) {
		return nil
	}
	var product model.MonitorProduct
	db.Where(model.MonitorProduct{Abbreviation: abbreviation}).First(&product)
	return &product

}
func (mpd *MonitorProductDao) SelectMonitorProductList() *[]model.MonitorProduct {
	var product = &[]model.MonitorProduct{}
	global.DB.Where("status = ?", "1").Find(product)
	return product
}
