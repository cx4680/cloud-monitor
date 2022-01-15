package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"gorm.io/gorm"
)

type AlarmInfoService struct {
	commonService.AbstractSyncServiceImpl
	AlarmInfoDao *commonDao.AlarmInfoDao
}

func NewAlarmInfoService() *AlarmInfoService {
	return &AlarmInfoService{AlarmInfoDao: commonDao.AlarmInfo}
}

func (s *AlarmInfoService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	list := param.([]commonModels.AlarmInfo)
	if len(list) > 0 {
		s.AlarmInfoDao.InsertBatch(db, list)
	}
	return jsonutil.ToString(list), nil
}
