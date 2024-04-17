package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"gorm.io/gorm"
	"strconv"
)

type ContactGroupRelService struct {
	service.AbstractSyncServiceImpl
	dao *dao.ContactGroupRelDao
}

func NewContactGroupRelService() *ContactGroupRelService {
	return &ContactGroupRelService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		dao:                     dao.ContactGroupRel,
	}
}

func (s *ContactGroupRelService) InsertRelByContact(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	relList, err := s.buildRelListByContact(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.InsertBatch(db, relList)
	return relList, nil
}

func (s *ContactGroupRelService) InsertRelByGroup(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	relList, err := s.buildRelListByGroup(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.InsertBatch(db, relList)
	return relList, nil
}

func (s *ContactGroupRelService) UpdateRelContact(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	relList, err := s.buildRelListByContact(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.UpdateByContact(db, relList, p)
	return relList, nil
}

func (s *ContactGroupRelService) UpdateRelGroup(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	relList, err := s.buildRelListByGroup(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.UpdateByGroup(db, relList, p)
	return relList, nil
}

func (s *ContactGroupRelService) DeleteRelByContact(db *gorm.DB, p *form.ContactParam) *model.ContactGroupRel {
	var contactGroupRel = &model.ContactGroupRel{TenantId: p.TenantId, ContactBizId: p.ContactBizId}
	s.dao.DeleteByContact(db, contactGroupRel)
	return contactGroupRel
}

func (s *ContactGroupRelService) DeleteRelByGroup(db *gorm.DB, p *form.ContactParam) *model.ContactGroupRel {
	var contactGroupRel = &model.ContactGroupRel{TenantId: p.TenantId, GroupBizId: p.GroupBizId}
	s.dao.DeleteByGroup(db, contactGroupRel)
	return contactGroupRel
}

func (s *ContactGroupRelService) buildRelListByContact(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	if len(p.GroupBizIdList) > constant.MaxContactGroup {
		return nil, errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
	}
	var list []*model.ContactGroupRel
	for _, groupBizId := range p.GroupBizIdList {
		contactGroupList := *s.dao.GetContactGroup(db, p.TenantId, groupBizId)
		if len(contactGroupList) == 0 {
			return nil, errors.NewBusinessError("该租户无此联系组")
		}
		if contactGroupList[0].ContactCount >= constant.MaxContactGroup {
			return nil, errors.NewBusinessError("有联系组已有" + strconv.Itoa(constant.MaxContactGroup) + "个联系人")
		}
		contactGroupRel := &model.ContactGroupRel{
			TenantId:     p.TenantId,
			ContactBizId: p.ContactBizId,
			GroupBizId:   groupBizId,
			CreateUser:   p.CreateUser,
		}
		list = append(list, contactGroupRel)
	}
	return list, nil
}

func (s *ContactGroupRelService) buildRelListByGroup(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	if len(p.ContactBizIdList) > constant.MaxContactGroup {
		return nil, errors.NewBusinessError("联系组限制添加" + strconv.Itoa(constant.MaxContactGroup) + "个联系人")
	}
	var list []*model.ContactGroupRel
	for _, contactBizId := range p.ContactBizIdList {
		if strutil.IsBlank(contactBizId) {
			continue
		}
		contactList := *s.dao.GetContact(db, p.TenantId, contactBizId, p.GroupBizId)
		if len(contactList) == 0 {
			return nil, errors.NewBusinessError("该租户无此联系人")
		}
		if contactList[0].GroupCount >= constant.MaxContactGroup {
			return nil, errors.NewBusinessError("有联系人已加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
		}
		contactGroupRel := &model.ContactGroupRel{
			TenantId:     p.TenantId,
			ContactBizId: contactBizId,
			GroupBizId:   p.GroupBizId,
			CreateUser:   p.CreateUser,
		}
		list = append(list, contactGroupRel)
	}
	return list, nil
}
