package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"gorm.io/gorm"
	"strconv"
)

type ContactGroupService struct {
	service.AbstractSyncServiceImpl
	dao                    *dao.ContactGroupDao
	contactGroupRelService *ContactGroupRelService
}

func NewContactGroupService(contactGroupRelService *ContactGroupRelService) *ContactGroupService {
	return &ContactGroupService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		dao:                     dao.ContactGroup,
		contactGroupRelService:  contactGroupRelService,
	}
}

func (s *ContactGroupService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(form.ContactParam)
	//每个联系人最多加入5个联系组
	if len(p.GroupBizIdList) >= constant.MaxContactGroup {
		return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
	}
	switch p.EventEum {
	case enum.InsertContactGroup:
		//参数校验
		if strutil.IsBlank(p.GroupName) {
			return "", errors.NewBusinessError("联系组名字不能为空")
		}
		if len(p.ContactBizIdList) == 0 {
			return "", errors.NewBusinessError("请至少选择一位联系人")
		}
		//联系组限制创建10个
		var groupCount int64
		global.DB.Model(&model.ContactGroup{}).Where("tenant_id = ?", p.TenantId).Count(&groupCount)
		if groupCount >= constant.MaxGroupNum {
			return "", errors.NewBusinessError("联系组限制创建" + strconv.Itoa(constant.MaxGroupNum) + "个")
		}
		//联系组名不可重复
		var count int64
		global.DB.Model(&model.ContactGroup{}).Where("tenant_id = ? AND name = ?", p.TenantId, p.GroupName).Count(&count)
		if count >= 1 {
			return "", errors.NewBusinessError("联系组名重复")
		}
		contactGroup, err := s.insertContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//保存联系人组关联
		p.GroupBizId = contactGroup.BizId
		if err := s.contactGroupRelService.PersistenceInner(db, s.contactGroupRelService, sys_rocketmq.ContactGroupTopic, p); err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.InsertContactGroup,
			Data:     contactGroup,
		}), nil
	case enum.UpdateContactGroup:
		//参数校验
		if strutil.IsBlank(p.GroupName) {
			return "", errors.NewBusinessError("联系组名字不能为空")
		}
		//联系组名不可重复
		var count int64
		global.DB.Model(&model.ContactGroup{}).Where("tenant_id = ? AND biz_id != ? AND name = ?", p.TenantId, p.GroupBizId, p.GroupName).Count(&count)
		if count >= 1 {
			return "", errors.NewBusinessError("联系组名重复")
		}
		contactGroup, err := s.updateContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		if err := s.contactGroupRelService.PersistenceInner(db, s.contactGroupRelService, sys_rocketmq.ContactGroupTopic, p); err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.UpdateContactGroup,
			Data:     contactGroup,
		}), nil
	case enum.DeleteContactGroup:
		contactGroup, err := s.deleteContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		if err := s.contactGroupRelService.PersistenceInner(db, s.contactGroupRelService, sys_rocketmq.ContactGroupTopic, p); err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.DeleteContactGroup,
			Data:     contactGroup,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *ContactGroupService) SelectContactGroup(param form.ContactParam) *form.ContactFormPage {
	return s.dao.SelectContactGroup(global.DB, param)
}

func (s *ContactGroupService) SelectAlertGroupContact(param form.ContactParam) *form.ContactFormPage {
	return s.dao.SelectGroupContact(global.DB, param)
}

func (s *ContactGroupService) insertContactGroup(db *gorm.DB, p form.ContactParam) (*model.ContactGroup, error) {
	currentTime := util.GetNow()
	contactGroup := &model.ContactGroup{
		BizId:       strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		TenantId:    p.TenantId,
		Name:        p.GroupName,
		Description: p.Description,
		CreateUser:  p.CreateUser,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
	}
	s.dao.Insert(db, contactGroup)
	return contactGroup, nil
}

func (s *ContactGroupService) updateContactGroup(db *gorm.DB, p form.ContactParam) (*model.ContactGroup, error) {
	var oldContactGroup model.ContactGroup
	global.DB.Where("tenant_id = ? AND biz_id = ?", p.TenantId, p.GroupBizId).First(&oldContactGroup)
	if oldContactGroup == (model.ContactGroup{}) {
		return nil, errors.NewBusinessError("该租户无此联系组")
	}
	currentTime := util.GetNow()
	var contactGroup = &model.ContactGroup{
		Id:          oldContactGroup.Id,
		BizId:       p.GroupBizId,
		TenantId:    p.TenantId,
		Name:        p.GroupName,
		Description: p.Description,
		UpdateTime:  currentTime,
		CreateTime:  oldContactGroup.CreateTime,
		CreateUser:  oldContactGroup.CreateUser,
	}
	s.dao.Update(db, contactGroup)
	return contactGroup, nil
}

func (s *ContactGroupService) deleteContactGroup(db *gorm.DB, p form.ContactParam) (*model.ContactGroup, error) {
	var contactGroup = &model.ContactGroup{
		BizId:    p.GroupBizId,
		TenantId: p.TenantId,
	}
	s.dao.Delete(db, contactGroup)
	return contactGroup, nil
}
