package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"gorm.io/gorm"
	"regexp"
	"strconv"
)

type ContactService struct {
	service.AbstractSyncServiceImpl
	dao                       *dao.ContactDao
	contactGroupService       *ContactGroupService
	contactInformationService *ContactInformationService
	contactGroupRelService    *ContactGroupRelService
}

func NewContactService(contactGroupService *ContactGroupService, contactInformationService *ContactInformationService, contactGroupRelService *ContactGroupRelService) *ContactService {
	return &ContactService{
		AbstractSyncServiceImpl:   service.AbstractSyncServiceImpl{},
		dao:                       dao.Contact,
		contactGroupService:       contactGroupService,
		contactInformationService: contactInformationService,
		contactGroupRelService:    contactGroupRelService,
	}
}

func (s *ContactService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(*form.ContactParam)
	//每个联系人最多加入5个联系组
	if len(p.GroupBizIdList) > constant.MaxContactGroup {
		return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
	}
	switch p.EventEum {
	case enum.InsertContact:
		//参数校验
		if strutil.IsBlank(p.ContactName) {
			return "", errors.NewBusinessError("联系人名字不能为空")
		}
		if !s.checkContactName(p.ContactName) {
			return "", errors.NewBusinessError("联系人名字格式错误")
		}
		contact, err := s.insertContact(db, p)
		if err != nil {
			return "", err
		}
		p.ContactBizId = contact.BizId
		//创建联系方式
		informationList, err := s.contactInformationService.InsertContactInformation(db, p)
		if err != nil {
			return "", err
		}
		//创建联系人组关联
		relList, err := s.contactGroupRelService.InsertContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.InsertContact,
			Data:     ContactMsg{Contact: contact, ContactInformationList: informationList, ContactGroupRelList: relList},
		}), nil
	case enum.UpdateContact:
		//参数校验
		if strutil.IsBlank(p.ContactName) {
			return "", errors.NewBusinessError("联系人名字不能为空")
		}
		if !s.checkContactName(p.ContactName) {
			return "", errors.NewBusinessError("联系人名字格式错误")
		}
		contact, err := s.updateContact(db, *p)
		if err != nil {
			return "", err
		}
		//更新联系方式
		informationList, err := s.contactInformationService.UpdateContactInformation(db, p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		relList, err := s.contactGroupRelService.UpdateContactGroupRel(db, p)
		if err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.UpdateContact,
			Data:     ContactMsg{Contact: contact, ContactInformationList: informationList, ContactGroupRelList: relList},
		}), nil
	case enum.DeleteContact:
		contact, err := s.deleteContact(db, *p)
		if err != nil {
			return "", err
		}
		//删除联系方式
		informationList := s.contactInformationService.DeleteContactInformation(db, p)
		//删除联系人组关联
		relList := s.contactGroupRelService.DeleteContactGroupRel(db, p)
		msg := form.MqMsg{
			EventEum: enum.DeleteContact,
			Data:     ContactMsg{Contact: contact, ContactInformation: informationList, ContactGroupRel: relList, Param: p},
		}
		return jsonutil.ToString(msg), nil
	case enum.ActivateContact:
		tenantId := s.dao.ActivateContact(db, p.ActiveCode)
		if strutil.IsBlank(p.ActiveCode) || strutil.IsBlank(tenantId) {
			return "", errors.NewBusinessError("无效激活码")
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.ActivateContact,
			Data:     p.ActiveCode,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *ContactService) SelectContact(param form.ContactParam) *form.ContactFormPage {
	return s.dao.Select(global.DB, param)
}

func (s *ContactService) insertContact(db *gorm.DB, p *form.ContactParam) (*model.Contact, error) {
	//每个账号限制创建100个联系人
	if s.dao.GetContactCount(p.TenantId) >= constant.MaxContactNum {
		return nil, errors.NewBusinessError("联系人限制创建" + strconv.Itoa(constant.MaxContactNum) + "个")
	}
	currentTime := util.GetNow()
	contact := &model.Contact{
		BizId:       strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		TenantId:    p.TenantId,
		Name:        p.ContactName,
		Description: p.Description,
		CreateUser:  p.CreateUser,
		State:       1,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
	}
	s.dao.Insert(db, contact)
	return contact, nil
}

func (s *ContactService) updateContact(db *gorm.DB, p form.ContactParam) (*model.Contact, error) {
	if !s.dao.CheckContact(p.TenantId, p.ContactBizId) {
		return nil, errors.NewBusinessError("该租户无此联系人")
	}
	if strutil.IsBlank(p.ContactBizId) {
		return nil, errors.NewBusinessError("联系人ID不能为空")
	}
	currentTime := util.GetNow()
	var contact = &model.Contact{
		BizId:       p.ContactBizId,
		TenantId:    p.TenantId,
		Name:        p.ContactName,
		State:       1,
		Description: p.Description,
		UpdateTime:  currentTime,
	}
	s.dao.Update(db, contact)
	return contact, nil
}

func (s *ContactService) deleteContact(db *gorm.DB, p form.ContactParam) (*model.Contact, error) {
	if !s.dao.CheckContact(p.TenantId, p.ContactBizId) {
		return nil, errors.NewBusinessError("该租户无此联系人")
	}
	if strutil.IsBlank(p.ContactBizId) {
		return nil, errors.NewBusinessError("联系人ID不能为空")
	}
	var contact = &model.Contact{
		BizId:    p.ContactBizId,
		TenantId: p.TenantId,
	}
	s.dao.Delete(db, contact)
	return contact, nil
}

func (s *ContactService) GetTenantId(activeCode string) string {
	return s.dao.GetTenantIdByActiveCode(activeCode)
}

func (s *ContactService) checkContactName(contactName string) bool {
	if regexp.MustCompile("^[a-zA-Z0-9_\u4e00-\u9fa5]{1,40}$").MatchString(contactName) {
		return true
	}
	return false
}

type ContactMsg struct {
	Param                  *form.ContactParam
	Contact                *model.Contact
	ContactInformation     *model.ContactInformation
	ContactGroupRel        *model.ContactGroupRel
	ContactInformationList []*model.ContactInformation
	ContactGroupRelList    []*model.ContactGroupRel
}
