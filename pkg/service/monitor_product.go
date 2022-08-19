package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/vo"
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

func (s *MonitorProductService) GetMonitorProduct() *[]model.MonitorProductDTO {
	return s.dao.GetMonitorProductDTO()
}

func (s *MonitorProductService) GetAllMonitorProduct() *[]model.MonitorProduct {
	return s.dao.GetAllMonitorProduct()
}

func (s *MonitorProductService) ChangeStatus(param form.MonitorProductParam) {
	s.dao.ChangeStatus(global.DB, param.ProductCodeList, param.Status)
}

func (s *MonitorProductService) GetMonitorProductPage(pageSize int, pageNum int) *vo.PageVO {
	var productList []model.MonitorProduct
	var total int64
	global.DB.Model(productList).Where("status = ?", "1").Count(&total)
	if total != 0 {
		offset := (pageNum - 1) * pageSize
		global.DB.Where("status = ?", "1").Order("sort ASC").Offset(offset).Limit(pageSize).Find(&productList)
	}
	return &vo.PageVO{
		Records: productList,
		Current: pageNum,
		Size:    pageSize,
		Total:   int(total),
	}
}
