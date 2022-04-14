package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	form2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	model2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"gorm.io/gorm"
	"regexp"
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
	p := param.(*form2.ContactParam)
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
		if !s.checkGroupName(p.GroupName) {
			return "", errors.NewBusinessError("联系组名字格式错误")
		}
		if len(p.ContactBizIdList) == 0 {
			return "", errors.NewBusinessError("请至少选择一位联系人")
		}
		//联系组限制创建10个
		if s.dao.GetGroupCount(p.TenantId) >= constant.MaxGroupNum {
			return "", errors.NewBusinessError("联系组限制创建" + strconv.Itoa(constant.MaxGroupNum) + "个")
		}
		//联系组名不可重复
		if s.dao.CheckGroupName(p.TenantId, p.GroupName, "") {
			return "", errors.NewBusinessError("联系组名重复")
		}
		contactGroup, err := s.insertContactGroup(db, *p)
		if err != nil {
			return "", err
		}
		p.GroupBizId = contactGroup.BizId
		//创建联系人组关联
		relList, err := s.contactGroupRelService.InsertContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return jsonutil.ToString(form2.MqMsg{
			EventEum: enum.InsertContactGroup,
			Data:     ContactGroupMsg{ContactGroup: contactGroup, ContactGroupRelList: relList},
		}), nil
	case enum.UpdateContactGroup:
		//参数校验
		if strutil.IsBlank(p.GroupName) {
			return "", errors.NewBusinessError("联系组名字不能为空")
		}
		if !s.checkGroupName(p.GroupName) {
			return "", errors.NewBusinessError("联系组名字格式错误")
		}
		//联系组名不可重复
		if s.dao.CheckGroupName(p.TenantId, p.GroupName, p.GroupBizId) {
			return "", errors.NewBusinessError("联系组名重复")
		}
		contactGroup, err := s.updateContactGroup(db, *p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		relList, err := s.contactGroupRelService.UpdateContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return jsonutil.ToString(form2.MqMsg{
			EventEum: enum.UpdateContactGroup,
			Data:     ContactGroupMsg{ContactGroup: contactGroup, ContactGroupRelList: relList},
		}), nil
	case enum.DeleteContactGroup:
		contactGroup, err := s.deleteContactGroup(db, *p)
		if err != nil {
			return "", err
		}
		//删除联系人组关联
		relList, err := s.contactGroupRelService.UpdateContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return jsonutil.ToString(form2.MqMsg{
			EventEum: enum.DeleteContactGroup,
			Data:     ContactGroupMsg{ContactGroup: contactGroup, ContactGroupRelList: relList},
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *ContactGroupService) SelectContactGroup(param form2.ContactParam) *form2.ContactFormPage {
	return s.dao.SelectContactGroup(global.DB, param)
}

func (s *ContactGroupService) SelectAlertGroupContact(param form2.ContactParam) *form2.ContactFormPage {
	return s.dao.SelectGroupContact(global.DB, param)
}

func (s *ContactGroupService) insertContactGroup(db *gorm.DB, p form2.ContactParam) (*model2.ContactGroup, error) {
	currentTime := util.GetNow()
	contactGroup := &model2.ContactGroup{
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

func (s *ContactGroupService) updateContactGroup(db *gorm.DB, p form2.ContactParam) (*model2.ContactGroup, error) {
	if strutil.IsBlank(p.GroupBizId) {
		return nil, errors.NewBusinessError("联系组ID不能为空")
	}
	var oldContactGroup = s.dao.GetGroup(p.TenantId, p.GroupBizId)
	if oldContactGroup == (model2.ContactGroup{}) {
		return nil, errors.NewBusinessError("该租户无此联系组")
	}
	currentTime := util.GetNow()
	var contactGroup = &model2.ContactGroup{
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

func (s *ContactGroupService) deleteContactGroup(db *gorm.DB, p form2.ContactParam) (*model2.ContactGroup, error) {
	//检验联系组是否存在
	if !s.dao.CheckGroupId(p.TenantId, p.GroupBizId) {
		return nil, errors.NewBusinessError("该租户无此联系组")
	}
	var contactGroup = &model2.ContactGroup{
		BizId:    p.GroupBizId,
		TenantId: p.TenantId,
	}
	s.dao.Delete(db, contactGroup)
	return contactGroup, nil
}

func (s *ContactGroupService) checkGroupName(groupName string) bool {
	if regexp.MustCompile("^[a-zA-Z0-9_\u4e00-\u9fa5]{1,40}$").MatchString(groupName) {
		return true
	}
	return false
}

type ContactGroupMsg struct {
	Param               *form2.ContactParam
	ContactGroup        *model2.ContactGroup
	ContactGroupRel     *model2.ContactGroupRel
	ContactGroupRelList []*model2.ContactGroupRel
}
