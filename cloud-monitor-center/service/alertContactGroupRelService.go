package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactGroupRelService struct {
	service.AbstractSyncServiceImpl
	dao *dao.AlertContactGroupRelDao
}

func NewAlertContactGroupRelService() *AlertContactGroupRelService {
	return &AlertContactGroupRelService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		dao:                     dao.AlertContactGroupRel,
	}
}

func (s *AlertContactGroupRelService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(forms.AlertContactParam)
	switch p.EventEum {
	case enums.InsertAlertContact:
		relList, err := s.insertAlertContactRel(db, p)
		if err != nil {
			return "", err
		}
		msg := forms.MqMsg{
			EventEum: enums.InsertAlertContactGroupRel,
			Data:     relList,
		}
		return tools.ToString(msg), nil
	case enums.UpdateAlertContact:
		List, err := s.updateAlertContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.UpdateAlertContactGroupRel,
			Data:     List,
		}), nil
	case enums.DeleteAlertContact:
		alertContactGroupRel, err := s.deleteAlertContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.DeleteAlertContactGroupRel,
			Data:     alertContactGroupRel,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")

	}
}

func (s *AlertContactGroupRelService) insertAlertContactRel(db *gorm.DB, p forms.AlertContactParam) ([]*models.AlertContactGroupRel, error) {
	list, err := s.getAlertContactGroupRelList(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.InsertBatch(db, list)
	return list, nil
}

func (s *AlertContactGroupRelService) updateAlertContactGroupRel(db *gorm.DB, p forms.AlertContactParam) ([]*models.AlertContactGroupRel, error) {
	list, err := s.getAlertContactGroupRelList(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.Update(db, list)
	return list, nil
}

func (s *AlertContactGroupRelService) deleteAlertContactGroupRel(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContactGroupRel, error) {
	if p.ContactId != "" {
		var alertContactGroupRel = &models.AlertContactGroupRel{
			ContactId: p.ContactId,
			TenantId:  p.TenantId,
		}
		s.dao.Delete(db, alertContactGroupRel)
		return alertContactGroupRel, nil
	} else {
		var alertContactGroupRel = &models.AlertContactGroupRel{
			GroupId:  p.GroupId,
			TenantId: p.TenantId,
		}
		s.dao.Delete(db, alertContactGroupRel)
		return alertContactGroupRel, nil
	}
}

func (s *AlertContactGroupRelService) getAlertContactGroupRelList(db *gorm.DB, p forms.AlertContactParam) ([]*models.AlertContactGroupRel, error) {
	var infoList []*models.AlertContactGroupRel
	var count int64
	for _, v := range p.GroupIdList {
		db.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ?", p.TenantId).Where("group_id = ?", v).Count(&count)
		if count >= constants.MaxContactNum {
			return nil, errors.NewBusinessError("每组联系人限制创建" + strconv.Itoa(constants.MaxContactNum) + "个")
		}
		alertContactGroup := &models.AlertContactGroupRel{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   p.TenantId,
			ContactId:  p.ContactId,
			GroupId:    v,
			CreateUser: p.CreateUser,
		}
		infoList = append(infoList, alertContactGroup)
	}
	return infoList, nil
}
