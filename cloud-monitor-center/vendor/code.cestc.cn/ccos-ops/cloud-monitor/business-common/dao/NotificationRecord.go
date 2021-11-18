package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"gorm.io/gorm"
	"time"
)

type NotificationRecordDao struct {
	db *gorm.DB
}

func NewNotificationRecordDao(db *gorm.DB) *NotificationRecordDao {
	return &NotificationRecordDao{db: db}
}

func (dao *NotificationRecordDao) InsertBatch(recordList []models.NotificationRecord) {
	dao.db.Create(&recordList)
}

func (dao *NotificationRecordDao) GetTenantPhoneCurrentMonthRecordNum(tenantId string) int {
	var count int64
	start, end := tools.GetMonthStartEnd(time.Now())
	dao.db.Debug().Model(&models.NotificationRecord{}).Where("sender_id=? and notification_type = ? and source_type != ? and create_time >= ? and create_time <= ?", tenantId, 1, dtos.SMS_LACK, start, end).Count(&count)
	return int(count)
}

func (dao *NotificationRecordDao) GetTenantSMSLackRecordNum(tenantId string) int {
	var count int64
	start, end := tools.GetMonthStartEnd(time.Now())
	dao.db.Debug().Model(&models.NotificationRecord{}).Where("sender_id=? and source_type = ? and create_time >= ? and create_time <= ?", tenantId, dtos.SMS_LACK, start, end).Count(&count)
	return int(count)
}
