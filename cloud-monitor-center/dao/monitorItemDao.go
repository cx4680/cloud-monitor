package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/models"
)

type MonitorItemDao struct {
}

var MonitorItem = new(MonitorItemDao)

func (mpd *MonitorItemDao) SelectMonitorItemsById(productId string) *[]models.MonitorItem {
	var product = &[]models.MonitorItem{}
	global.DB.Where("status = ?", "1").Where("is_display = ?", "1").Where("product_id = ?", productId).Find(product)
	return product
}
