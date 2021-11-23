package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type AlertContactInformationDao struct {
}

var AlertContactInformation = new(AlertContactInformationDao)

func (acid *AlertContactInformationDao) Insert(db *gorm.DB, entity *models.AlertContactInformation) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	entity.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (acid *AlertContactInformationDao) InsertBatch(db *gorm.DB, list []*models.AlertContactInformation) {
	for _, information := range list {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		information.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
		information.CreateTime = currentTime
		information.UpdateTime = currentTime
	}
	db.Create(list)
}
