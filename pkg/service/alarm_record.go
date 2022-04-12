package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"context"
	"gorm.io/gorm"
)

type AlarmRecordService struct {
	AlarmRecordDao *commonDao.AlarmRecordDao
	AlarmInfoDao   *commonDao.AlarmInfoDao
}

func NewAlarmRecordService() *AlarmRecordService {
	return &AlarmRecordService{AlarmRecordDao: commonDao.AlarmRecord, AlarmInfoDao: commonDao.AlarmInfo}
}

func (s *AlarmRecordService) InsertAndHandler(ctx *context.Context, recordList []commonModels.AlarmRecord, infoList []commonModels.AlarmInfo, eventList []interface{}) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if len(recordList) > 0 {
			s.AlarmRecordDao.InsertBatch(tx, recordList)
		}
		if len(infoList) > 0 {
			s.AlarmInfoDao.InsertBatch(tx, infoList)
		}
		//告警处置
		if !AlarmHandlerQueue.PushFrontBatch(eventList) {
			logger.Logger().Error("requestId=", util.GetRequestId(ctx), ", add to alarm handler fail, handlerMap=", jsonutil.ToString(eventList))
		}
		return nil
	})
}
