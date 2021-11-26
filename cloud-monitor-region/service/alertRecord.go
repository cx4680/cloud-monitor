package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"gorm.io/gorm"
)

type AlertRecordService struct {
	commonService.AbstractSyncServiceImpl
	AlertRecordDao *commonDao.AlertRecordDao
}

func NewAlertRecordService() *AlertRecordService {
	return &AlertRecordService{AlertRecordDao: commonDao.AlertRecord}
}

func (s *AlertRecordService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	list := param.([]commonModels.AlertRecord)
	if len(list) > 0 {
		s.AlertRecordDao.InsertBatch(db, list)
	}
	return tools.ToString(list), nil
}
