package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	commonForm "code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"gorm.io/gorm"
)

type MonitorItemService struct {
	dao *dao.MonitorItemDao
	service.AbstractSyncServiceImpl
}

func NewMonitorItemService(dao *dao.MonitorItemDao) *MonitorItemService {
	return &MonitorItemService{
		dao:                     dao,
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{}}
}

func (s *MonitorItemService) GetMonitorItem(param form.MonitorItemParam) []model.MonitorItem {
	return s.dao.GetMonitorItem(param.ProductBizId, param.OsType, param.Display)
}

func (s *MonitorItemService) ChangeDisplay(db *gorm.DB, param form.MonitorItemParam) {
	s.dao.ChangeDisplay(db, param.ProductBizId, param.Display, param.BizIdList)
}

func (s *MonitorItemService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(commonForm.MonitorItemParam)
	switch p.EventEum {
	case enum.ChangeMonitorItemDisplay:
		s.ChangeDisplay(db, p)
		msg := commonForm.MqMsg{
			EventEum: enum.ChangeMonitorItemDisplay,
			Data:     param,
		}
		return jsonutil.ToString(msg), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *MonitorItemService) GetMonitorItemPage(pageSize int, pageNum int, productAbbr string) *vo.PageVO {
	var monitorItemList []model.MonitorItem
	sql := "select item.* from t_monitor_item item  ,t_monitor_product product  where item.product_biz_id=product.biz_id  and product.abbreviation=?"
	paginate := util.Paginate(pageSize, pageNum, sql, []interface{}{productAbbr}, &monitorItemList)
	return paginate
}
