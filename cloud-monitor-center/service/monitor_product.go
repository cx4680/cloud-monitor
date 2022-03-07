package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	commonForm "code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"gorm.io/gorm"
)

type MonitorProductService struct {
	dao *dao.MonitorProductDao
	service.AbstractSyncServiceImpl
}

func NewMonitorProductService(dao *dao.MonitorProductDao) *MonitorProductService {
	return &MonitorProductService{
		dao:                     dao,
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{}}
}

func (s *MonitorProductService) GetMonitorProduct() *[]model.MonitorProduct {
	return s.dao.GetMonitorProduct()
}

func (s *MonitorProductService) GetAllMonitorProduct() *[]model.MonitorProduct {
	return s.dao.GetAllMonitorProduct()
}

func (s *MonitorProductService) ChangeStatus(param commonForm.MonitorProductParam) {
	s.dao.ChangeStatus(param.BizIdList, param.Status)
}

func (s *MonitorProductService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(commonForm.MonitorProductParam)
	switch p.EventEum {
	case enum.ChangeMonitorProductStatus:
		s.ChangeStatus(p)
		msg := commonForm.MqMsg{
			EventEum: enum.ChangeMonitorProductStatus,
			Data:     param,
		}
		return jsonutil.ToString(msg), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}
