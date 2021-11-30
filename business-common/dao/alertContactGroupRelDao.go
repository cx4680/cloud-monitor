package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"gorm.io/gorm"
)

type AlertContactGroupRelDao struct {
}

var AlertContactGroupRel = new(AlertContactGroupRelDao)

func (acid *AlertContactGroupRelDao) Insert(db *gorm.DB, entity *models.AlertContactGroupRel) {
	currentTime := tools.GetNowStr()
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (acid *AlertContactGroupRelDao) InsertBatch(db *gorm.DB, list []*models.AlertContactGroupRel) {
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

func (acid *AlertContactGroupRelDao) Update(db *gorm.DB, list []*models.AlertContactGroupRel, entity *models.AlertContactGroupRel) {
	acid.Delete(db, entity)
	acid.InsertBatch(db, list)
}

func (acid *AlertContactGroupRelDao) Delete(db *gorm.DB, entity *models.AlertContactGroupRel) {
	db.Delete(entity)
}
