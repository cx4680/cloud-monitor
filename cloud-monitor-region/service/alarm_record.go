package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"gorm.io/gorm"
)

type AlarmRecordService struct {
	commonService.AbstractSyncServiceImpl
	AlarmRecordDao *commonDao.AlarmRecordDao
	AlarmInfoSvc   *AlarmInfoService
}

func NewAlarmRecordService(AlarmInfoSvc *AlarmInfoService) *AlarmRecordService {
	return &AlarmRecordService{AlarmRecordDao: commonDao.AlarmRecord, AlarmInfoSvc: AlarmInfoSvc}
}

func (s *AlarmRecordService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	if p, ok := param.(AlarmAddParam); ok {
		if len(p.InfoList) > 0 {
			if err := s.AlarmInfoSvc.PersistenceInner(db, s.AlarmInfoSvc, sys_rocketmq.AlarmInfoTopic, p.InfoList); err != nil {
				return "", err
			}
		}
		if len(p.RecordList) > 0 {
			s.AlarmRecordDao.InsertBatch(db, p.RecordList)
			return jsonutil.ToString(p.RecordList), nil
		}
	} else {
		logger.Logger().Error("Alarm Record persistence param type error, ", jsonutil.ToString(param))
	}
	logger.Logger().Error("Alarm Record persistence fail, ", jsonutil.ToString(param))
	return "", nil
}
