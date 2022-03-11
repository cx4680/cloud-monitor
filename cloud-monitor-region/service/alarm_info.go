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
		//擦除自增Id，解决同步到中心化Id冲突问题
		for i, _ := range list {
			list[i].Id = 0
		}
	}
	return jsonutil.ToString(list), nil
}
