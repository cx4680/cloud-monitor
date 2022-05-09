package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/vo"
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

func (s *MonitorProductService) GetMonitorProduct() *[]model.MonitorProductDTO {
	return s.dao.GetMonitorProductDTO()
}

func (s *MonitorProductService) GetAllMonitorProduct() *[]model.MonitorProduct {
	return s.dao.GetAllMonitorProduct()
}

func (s *MonitorProductService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(form.MonitorProductParam)
	switch p.EventEum {
	case enum.ChangeMonitorProductStatus:
		s.dao.ChangeStatus(db, p.BizIdList, p.Status)
		msg := form.MqMsg{
			EventEum: enum.ChangeMonitorProductStatus,
			Data:     param,
		}
		return jsonutil.ToString(msg), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
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
