package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"gorm.io/gorm"
)

type AlertContactGroupRelDao struct {
}

var AlertContactGroupRel = new(AlertContactGroupRelDao)

func (d *AlertContactGroupRelDao) Insert(db *gorm.DB, entity *models.AlertContactGroupRel) {
	currentTime := tools.GetNowStr()
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (d *AlertContactGroupRelDao) InsertBatch(db *gorm.DB, list []*models.AlertContactGroupRel) {
	if len(list) == 0 {
		return
	}
	currentTime := tools.GetNowStr()
	for _, information := range list {
		information.CreateTime = currentTime
		information.UpdateTime = currentTime
	}
	db.Create(list)
}

func (d *AlertContactGroupRelDao) Update(db *gorm.DB, list []*models.AlertContactGroupRel, entity *models.AlertContactGroupRel) {
	d.Delete(db, entity)
	d.InsertBatch(db, list)
}

func (d *AlertContactGroupRelDao) Delete(db *gorm.DB, entity *models.AlertContactGroupRel) {
	if entity.ContactId != "" {
		db.Where("tenant_id = ? AND contact_id = ?", entity.TenantId, entity.ContactId).Delete(models.AlertContactGroupRel{})
	} else {
		db.Where("tenant_id = ? AND group_id = ?", entity.TenantId, entity.GroupId).Delete(models.AlertContactGroupRel{})
	}
}
