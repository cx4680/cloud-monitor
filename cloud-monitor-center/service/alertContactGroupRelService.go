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
	case enums.InsertAlertContact, enums.InsertAlertContactGroup:
		relList, err := s.insertAlertContactRel(db, p)
		if err != nil {
			return "", err
		}
		msg := forms.MqMsg{
			EventEum: enums.InsertAlertContactGroupRel,
			Data:     relList,
		}
		return tools.ToString(msg), nil
	case enums.UpdateAlertContact, enums.UpdateAlertContactGroup:
		List, err := s.updateAlertContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.UpdateAlertContactGroupRel,
			Data:     List,
		}), nil
	case enums.DeleteAlertContact, enums.DeleteAlertContactGroup:
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
	var list []*models.AlertContactGroupRel
	var count int64
	if len(p.ContactIdList) > 0 {
		if len(p.ContactIdList) > constants.MaxContactNum {
			return nil, errors.NewBusinessError("联系组限制添加" + strconv.Itoa(constants.MaxContactNum) + "个联系人")
		}
		for _, contactId := range p.ContactIdList {
			db.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ? AND contact_id = ?", p.TenantId, contactId).Count(&count)
			if count >= constants.MaxContactNum {
				return nil, errors.NewBusinessError("有联系人已加入" + strconv.Itoa(constants.MaxContactGroup) + "个联系组")
			}
			alertContactGroupRel := &models.AlertContactGroupRel{
				Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
				TenantId:   p.TenantId,
				ContactId:  contactId,
				GroupId:    p.GroupId,
				CreateUser: p.CreateUser,
			}
			list = append(list, alertContactGroupRel)
		}
	} else if len(p.GroupIdList) > 0 {
		if len(p.GroupIdList) > constants.MaxContactGroup {
			return nil, errors.NewBusinessError("联系人限制加入" + strconv.Itoa(constants.MaxContactGroup) + "个联系组")
		}
		for _, groupId := range p.GroupIdList {
			db.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ? AND group_id = ?", p.TenantId, groupId).Count(&count)
			if count >= constants.MaxContactNum {
				return nil, errors.NewBusinessError("有联系组已有" + strconv.Itoa(constants.MaxContactNum) + "个联系人")
			}
			alertContactGroupRel := &models.AlertContactGroupRel{
				Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
				TenantId:   p.TenantId,
				ContactId:  p.ContactId,
				GroupId:    groupId,
				CreateUser: p.CreateUser,
			}
			list = append(list, alertContactGroupRel)
		}
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
	var alertContactGroupRel = &models.AlertContactGroupRel{}
	if p.ContactId != "" {
		alertContactGroupRel.TenantId = p.TenantId
		alertContactGroupRel.ContactId = p.ContactId
	} else {
		alertContactGroupRel.TenantId = p.TenantId
		alertContactGroupRel.GroupId = p.GroupId
	}
	s.dao.Delete(db, alertContactGroupRel)
	return alertContactGroupRel, nil
}

func (s *AlertContactGroupRelService) getAlertContactGroupRelList(db *gorm.DB, p forms.AlertContactParam) ([]*models.AlertContactGroupRel, error) {
	var list []*models.AlertContactGroupRel
	var count int64
	if len(p.ContactIdList) > 0 {
		for _, v := range p.ContactIdList {
			alertContactGroupRel := &models.AlertContactGroupRel{
				Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
				TenantId:   p.TenantId,
				ContactId:  v,
				GroupId:    p.GroupId,
				CreateUser: p.CreateUser,
			}
			list = append(list, alertContactGroupRel)
		}
	} else {
		for _, v := range p.GroupIdList {
			db.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ? AND contact_id = ?", p.TenantId, p.ContactId).Count(&count)
			if count >= constants.MaxContactNum {
				return nil, errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constants.MaxContactGroup) + "个联系组")
			}
			alertContactGroupRel := &models.AlertContactGroupRel{
				Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
				TenantId:   p.TenantId,
				ContactId:  p.ContactId,
				GroupId:    v,
				CreateUser: p.CreateUser,
			}
			list = append(list, alertContactGroupRel)
		}
	}
	return list, nil
}
