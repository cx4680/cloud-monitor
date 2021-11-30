package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactInformationDao struct {
}

var AlertContactInformation = new(AlertContactInformationDao)

func (acid *AlertContactInformationDao) Insert(db *gorm.DB, entity *models.AlertContactInformation) {
	currentTime := tools.GetNowStr()
	entity.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (acid *AlertContactInformationDao) InsertBatch(db *gorm.DB, list []*models.AlertContactInformation) {
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

func (acid *AlertContactInformationDao) Update(db *gorm.DB, list []*models.AlertContactInformation, entity *models.AlertContactInformation) {
	acid.Delete(db, entity)
	acid.InsertBatch(db, list)
}

func (acid *AlertContactInformationDao) Delete(db *gorm.DB, entity *models.AlertContactInformation) {
	db.Delete(entity)
}
