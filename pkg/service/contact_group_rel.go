package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
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

func (s *ContactGroupRelService) InsertContactGroupRel(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	list, err := s.buildRelList(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.InsertBatch(db, list)
	return list, nil
}

func (s *ContactGroupRelService) UpdateContactGroupRel(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	list, err := s.buildRelList(db, p)
	if err != nil {
		return nil, err
	}
	s.dao.Update(db, list, p)
	return list, nil
}

func (s *ContactGroupRelService) DeleteContactGroupRel(db *gorm.DB, p *form.ContactParam) *model.ContactGroupRel {
	var contactGroupRel = &model.ContactGroupRel{}
	if strutil.IsNotBlank(p.ContactBizId) {
		contactGroupRel.TenantId = p.TenantId
		contactGroupRel.ContactBizId = p.ContactBizId
	} else {
		contactGroupRel.TenantId = p.TenantId
		contactGroupRel.GroupBizId = p.GroupBizId
	}
	s.dao.Delete(db, contactGroupRel)
	return contactGroupRel
}

//构建组关联关系
func (s *ContactGroupRelService) buildRelList(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	var relList []*model.ContactGroupRel
	var err error
	if len(p.ContactBizIdList) > 0 {
		relList, err = s.buildContactRelList(db, p)
	} else if len(p.GroupBizIdList) > 0 {
		relList, err = s.buildGroupRelList(db, p)
	} else {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return relList, nil
}

func (s *ContactGroupRelService) buildContactRelList(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	if len(p.ContactBizIdList) > constant.MaxContactNum {
		return nil, errors.NewBusinessError("联系组限制添加" + strconv.Itoa(constant.MaxContactNum) + "个联系人")
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

func (s *ContactGroupRelService) buildGroupRelList(db *gorm.DB, p *form.ContactParam) ([]*model.ContactGroupRel, error) {
	if len(p.GroupBizIdList) > constant.MaxContactGroup {
		return nil, errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
	}
	var list []*model.ContactGroupRel
	for _, groupBizId := range p.GroupBizIdList {
		contactGroupList := *s.dao.GetContactGroup(db, p.TenantId, groupBizId)
		if len(contactGroupList) == 0 {
			return nil, errors.NewBusinessError("该租户无此联系组")
		}
		if contactGroupList[0].ContactCount >= constant.MaxContactNum {
			return nil, errors.NewBusinessError("有联系组已有" + strconv.Itoa(constant.MaxContactNum) + "个联系人")
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
