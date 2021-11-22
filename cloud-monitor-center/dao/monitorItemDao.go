package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
	"gorm.io/gorm"
)

type MonitorItemDao struct {
	db *gorm.DB
}

func NewMonitorItemDao() *MonitorItemDao {
	return &MonitorItemDao{db: database.GetDb()}
}

func (mpd *MonitorItemDao) SelectMonitorItemsById(productId string) *[]models.MonitorItem {
	var product = &[]models.MonitorItem{}
	mpd.db.Where("status = ?", "1").Where("is_display = ?", "1").Where("product_id = ?", productId).Find(product)
	return product
}
