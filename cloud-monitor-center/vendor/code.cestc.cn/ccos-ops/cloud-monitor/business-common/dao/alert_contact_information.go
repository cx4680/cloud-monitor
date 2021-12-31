package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactInformationDao struct {
}

var AlertContactInformation = new(AlertContactInformationDao)

func (d *AlertContactInformationDao) Insert(db *gorm.DB, entity *model.AlertContactInformation) {
	currentTime := util.GetNowStr()
	entity.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (d *AlertContactInformationDao) InsertBatch(db *gorm.DB, list []*model.AlertContactInformation) {
	if len(list) == 0 {
		return
	}
	currentTime := util.GetNowStr()
	for _, v := range list {
		v.CreateTime = currentTime
		v.UpdateTime = currentTime
	}
	db.Create(list)
}

func (d *AlertContactInformationDao) Update(db *gorm.DB, list []*model.AlertContactInformation) {
	if len(list) == 0 {
		return
	}
	d.Delete(db, list[0])
	d.InsertBatch(db, list)
}

func (d *AlertContactInformationDao) Delete(db *gorm.DB, entity *model.AlertContactInformation) {
	db.Where("tenant_id = ? AND contact_id = ?", entity.TenantId, entity.ContactId).Delete(model.AlertContactInformation{})
}
