package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
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
	p := param.(form.AlertContactParam)
	switch p.EventEum {
	case enum.InsertAlertContact, enum.InsertAlertContactGroup:
		relList, err := s.insertAlertContactRel(db, p)
		if err != nil {
			return "", err
		}
		msg := form.MqMsg{
			EventEum: enum.InsertAlertContactGroupRel,
			Data:     relList,
		}
		return jsonutil.ToString(msg), nil
	case enum.UpdateAlertContact, enum.UpdateAlertContactGroup:
		List, err := s.updateAlertContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		var date = model.UpdateAlertContactGroupRel{
			RelList: List,
			Param:   p,
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.UpdateAlertContactGroupRel,
			Data:     date,
		}), nil
	case enum.DeleteAlertContact, enum.DeleteAlertContactGroup:
		alertContactGroupRel, err := s.deleteAlertContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.DeleteAlertContactGroupRel,
			Data:     alertContactGroupRel,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *AlertContactGroupRelService) insertAlertContactRel(db *gorm.DB, p form.AlertContactParam) ([]*model.AlertContactGroupRel, error) {
	list, err := s.buildRelList(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.InsertBatch(db, list)
	return list, nil
}

func (s *AlertContactGroupRelService) updateAlertContactGroupRel(db *gorm.DB, p form.AlertContactParam) ([]*model.AlertContactGroupRel, error) {
	list, err := s.buildRelList(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.Update(db, list, p)
	return list, nil
}

func (s *AlertContactGroupRelService) deleteAlertContactGroupRel(db *gorm.DB, p form.AlertContactParam) (*model.AlertContactGroupRel, error) {
	var alertContactGroupRel = &model.AlertContactGroupRel{}
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

//构建组关联关系
func (s *AlertContactGroupRelService) buildRelList(db *gorm.DB, p form.AlertContactParam) ([]*model.AlertContactGroupRel, error) {
	var relList []*model.AlertContactGroupRel
	var err error
	if len(p.ContactIdList) > 0 {
		relList, err = s.buildContactRelList(db, p)
	} else if len(p.GroupIdList) > 0 {
		relList, err = s.buildGroupRelList(db, p)
	} else {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return relList, nil
}

func (s *AlertContactGroupRelService) buildContactRelList(db *gorm.DB, p form.AlertContactParam) ([]*model.AlertContactGroupRel, error) {
	if len(p.ContactIdList) > constant.MaxContactNum {
		return nil, errors.NewBusinessError("联系组限制添加" + strconv.Itoa(constant.MaxContactNum) + "个联系人")
	}
	var list []*model.AlertContactGroupRel
	for _, contactId := range p.ContactIdList {
		if strutil.IsBlank(contactId) {
			continue
		}
		alertContactList := *s.dao.GetAlertContact(db, p.TenantId, contactId, p.GroupId)
		if len(alertContactList) == 0 {
			return nil, errors.NewBusinessError("该租户无此联系人")
		}
		if alertContactList[0].GroupCount >= constant.MaxContactGroup {
			return nil, errors.NewBusinessError("有联系人已加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
		}
		alertContactGroupRel := &model.AlertContactGroupRel{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   p.TenantId,
			ContactId:  contactId,
			GroupId:    p.GroupId,
			CreateUser: p.CreateUser,
		}
		list = append(list, alertContactGroupRel)
	}
	return list, nil
}

func (s *AlertContactGroupRelService) buildGroupRelList(db *gorm.DB, p form.AlertContactParam) ([]*model.AlertContactGroupRel, error) {
	if len(p.GroupIdList) > constant.MaxContactGroup {
		return nil, errors.NewBusinessError("联系人限制加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
	}
	var list []*model.AlertContactGroupRel
	for _, groupId := range p.GroupIdList {
		contactGroupList := *s.dao.GetAlertContactGroup(db, p.TenantId, groupId)
		if len(contactGroupList) == 0 {
			return nil, errors.NewBusinessError("该租户无此联系组")
		}
		if contactGroupList[0].ContactCount >= constant.MaxContactNum {
			return nil, errors.NewBusinessError("有联系组已有" + strconv.Itoa(constant.MaxContactNum) + "个联系人")
		}
		alertContactGroupRel := &model.AlertContactGroupRel{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   p.TenantId,
			ContactId:  p.ContactId,
			GroupId:    groupId,
			CreateUser: p.CreateUser,
		}
		list = append(list, alertContactGroupRel)
	}
	return list, nil
}
