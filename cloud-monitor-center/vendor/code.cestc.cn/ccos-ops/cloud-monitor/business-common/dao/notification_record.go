package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type NotificationRecordDao struct {
}

var NotificationRecord = new(NotificationRecordDao)

func (dao *NotificationRecordDao) InsertBatch(db *gorm.DB, recordList []model.NotificationRecord) {
	if len(recordList) > 0 {
		db.Create(&recordList)
	}
}

func (dao *NotificationRecordDao) Insert(db *gorm.DB, record model.NotificationRecord) {
	if strutil.IsBlank(record.Id) {
		record.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	}
	if strutil.IsBlank(record.CreateTime) {
		record.CreateTime = util.GetNowStr()
	}
	db.Create(&record)
}

func (dao *NotificationRecordDao) GetTenantPhoneCurrentMonthRecordNum(tenantId string) int {
	var count int64
	start, end := util.GetMonthStartEnd(time.Now())
	global.DB.Debug().Model(&model.NotificationRecord{}).Where("sender_id=? and notification_type = ? and source_type != ? and create_time >= ? and create_time <= ?", tenantId, 1, dto.SMS_LACK, start, end).Count(&count)
	return int(count)
}

func (dao *NotificationRecordDao) GetTenantSMSLackRecordNum(tenantId string) int {
	var count int64
	start, end := util.GetMonthStartEnd(time.Now())
	global.DB.Debug().Model(&model.NotificationRecord{}).Where("sender_id=? and source_type = ? and create_time >= ? and create_time <= ?", tenantId, dto.SMS_LACK, start, end).Count(&count)
	return int(count)
}
