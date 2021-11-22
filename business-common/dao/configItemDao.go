package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
)

func GetConfigItem(code interface{}, pid string, data string) *models.ConfigItem {
	item := models.ConfigItem{}
	db := global.DB
	if code != nil {
		db = db.Where("code", code)
	}
	if len(pid) > 0 {
		db = db.Where("pid", pid)
	}
	if len(data) > 0 {
		db = db.Where("data", data)
	}
	db.Find(&item)
	return &item
}
