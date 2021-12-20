package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactService struct {
	service.AbstractSyncServiceImpl
	dao                            *dao.AlertContactDao
	alertContactGroupService       *AlertContactGroupService
	alertContactInformationService *AlertContactInformationService
	alertContactGroupRelService    *AlertContactGroupRelService
}

func NewAlertContactService(alertContactGroupService *AlertContactGroupService, alertContactInformationService *AlertContactInformationService, alertContactGroupRelService *AlertContactGroupRelService) *AlertContactService {
	return &AlertContactService{
		AbstractSyncServiceImpl:        service.AbstractSyncServiceImpl{},
		dao:                            dao.AlertContact,
		alertContactGroupService:       alertContactGroupService,
		alertContactInformationService: alertContactInformationService,
		alertContactGroupRelService:    alertContactGroupRelService,
	}
}

func (s *AlertContactService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(forms.AlertContactParam)
	//每个联系人最多加入5个联系组
	if len(p.GroupIdList) >= constants.MaxContactGroup {
		return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constants.MaxContactGroup) + "个联系组")
	}
	switch p.EventEum {
	case enums.InsertAlertContact:
		//参数校验
		if p.ContactName == "" {
			return "", errors.NewBusinessError("联系人名字不能为空")
		}
		alertContact, err := s.insertAlertContact(db, p)
		if err != nil {
			return "", err
		}
		p.ContactId = alertContact.Id

		//保存 联系方式 和 联系人组关联
		if err := s.persistenceInner(db, p); err != nil {
			return "", err
		}
		msg := forms.MqMsg{
			EventEum: enums.InsertAlertContact,
			Data:     alertContact,
		}
		return tools.ToString(msg), nil
	case enums.UpdateAlertContact:
		//参数校验
		if p.ContactName == "" {
			return "", errors.NewBusinessError("联系人名字不能为空")
		}
		alertContact, err := s.updateAlertContact(db, p)
		if err != nil {
			return "", err
		}
		//更新 联系方式 和 联系人组关联
		if err := s.persistenceInner(db, p); err != nil {
			return "", err
		}
		msg := forms.MqMsg{
			EventEum: enums.UpdateAlertContact,
			Data:     alertContact,
		}
		return tools.ToString(msg), nil
	case enums.DeleteAlertContact:
		alertContact, err := s.deleteAlertContact(db, p)
		if err != nil {
			return "", err
		}
		//删除 联系方式 和 联系人组关联
		if err := s.persistenceInner(db, p); err != nil {
			return "", err
		}
		msg := forms.MqMsg{
			EventEum: enums.DeleteAlertContact,
			Data:     alertContact,
		}
		return tools.ToString(msg), nil
	case enums.CertifyAlertContact:
		s.dao.CertifyAlertContact(db, p.ActiveCode)
		msg := forms.MqMsg{
			EventEum: enums.CertifyAlertContact,
			Data:     p.ActiveCode,
		}
		return tools.ToString(msg), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *AlertContactService) SelectAlertContact(param forms.AlertContactParam) *forms.AlertContactFormPage {
	return s.dao.Select(global.DB, param)
}

func (s *AlertContactService) insertAlertContact(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContact, error) {
	//每个账号限制创建100个联系人
	var count int64
	global.DB.Model(&models.AlertContact{}).Where("tenant_id = ?", p.TenantId).Count(&count)
	if count >= constants.MaxContactNum {
		return nil, errors.NewBusinessError("联系人限制创建" + strconv.Itoa(constants.MaxContactNum) + "个")
	}
	currentTime := tools.GetNowStr()
	alertContact := &models.AlertContact{
		Id:          strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		TenantId:    p.TenantId,
		Name:        p.ContactName,
		Description: p.Description,
		CreateUser:  p.CreateUser,
		Status:      1,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
	}
	s.dao.Insert(db, alertContact)
	return alertContact, nil
}

func (s *AlertContactService) updateAlertContact(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContact, error) {
	//校验联系人是否为该租户所有
	var check int64
	db.Model(&models.AlertContact{}).Where("tenant_id = ? AND id = ?", p.TenantId, p.ContactId).Count(&check)
	if check == 0 {
		return nil, errors.NewBusinessError("该租户无此联系人")
	}

	currentTime := tools.GetNowStr()
	var alertContact = &models.AlertContact{
		Id:          p.ContactId,
		TenantId:    p.TenantId,
		Name:        p.ContactName,
		Status:      1,
		Description: p.Description,
		UpdateTime:  currentTime,
	}
	s.dao.Update(db, alertContact)
	return alertContact, nil
}

func (s *AlertContactService) deleteAlertContact(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContact, error) {
	//校验联系人是否为该租户所有
	var check int64
	db.Model(&models.AlertContact{}).Where("tenant_id = ? AND contact_id = ?", p.TenantId, p.ContactId).Count(&check)
	if p.ContactId == "" {
		return nil, errors.NewBusinessError("联系人ID不能为空")
	}
	var alertContact = &models.AlertContact{
		Id:       p.ContactId,
		TenantId: p.TenantId,
	}
	s.dao.Delete(db, alertContact)
	return alertContact, nil
}

func (s *AlertContactService) persistenceInner(db *gorm.DB, p forms.AlertContactParam) error {
	//联系人组关联
	if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sysRocketMq.AlertContactGroupTopic, p); err != nil {
		return err
	}
	//联系方式
	if err := s.alertContactInformationService.PersistenceInner(db, s.alertContactInformationService, sysRocketMq.AlertContactTopic, p); err != nil {
		return err
	}
	return nil
}

func (s *AlertContactService) GetTenantId(activeCode string) string {
	var model = &models.AlertContactInformation{}
	global.DB.Where("active_code = ?", activeCode).First(model)
	return model.TenantId
}
