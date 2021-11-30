package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
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
	switch p.EventEum {
	case enums.InsertAlertContact:
		alertContact, err := s.insertAlertContact(db, p)
		if err != nil {
			return "", err
		}
		////保存联系方式
		//if err := s.alertContactInformationService.PersistenceInner(db, s.alertContactInformationService, sysRocketMq.AlertContactTopic, p); err != nil {
		//	return "", err
		//}
		////保存联系人组关联
		//if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sysRocketMq.AlertContactTopic, p); err != nil {
		//	return "", err
		//}
		msg := forms.MqMsg{
			EventEum: enums.InsertAlertContact,
			Data:     alertContact,
		}
		return tools.ToString(msg), nil
	case enums.UpdateAlertContact:
		if p.ContactName == "" {
			return "", errors.NewBusinessError("联系人名字不能为空")
		}
		//每个联系人最多加入5个联系组
		if len(p.GroupIdList) >= constants.MaxContactGroup {
			return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constants.MaxContactGroup) + "个联系组")
		}
		var alertContact = &models.AlertContact{
			Id:          p.ContactId,
			TenantId:    p.TenantId,
			Name:        p.ContactName,
			Status:      1,
			Description: p.Description,
		}
		s.dao.Update(db, alertContact)

		//更新联系方式
		if err := s.alertContactInformationService.PersistenceInner(db, s.alertContactInformationService, sysRocketMq.AlertContactTopic, param); err != nil {
			return "", errors.NewBusinessError(err.Error())
		}
		//更新联系人组关联
		if err := s.alertContactGroupService.PersistenceInner(db, s.alertContactGroupService, sysRocketMq.AlertContactTopic, param); err != nil {
			return "", errors.NewBusinessError(err.Error())
		}

		msg := forms.MqMsg{
			EventEum: enums.UpdateAlertContact,
			Data:     alertContact,
		}
		return tools.ToString(msg), nil
	case enums.DeleteAlertContact:
		if p.ContactId == "" {
			return "", errors.NewBusinessError("联系人ID不能为空")
		}
		var alertContact = &models.AlertContact{
			Id:       p.ContactId,
			TenantId: p.TenantId,
		}
		s.dao.Delete(db, alertContact)
		msg := forms.MqMsg{
			EventEum: enums.DeleteAlertContact,
			Data:     alertContact,
		}
		return tools.ToString(msg), nil
	case enums.CertifyAlertContact:
		s.dao.CertifyAlertContact(p.ActiveCode)
		msg := forms.MqMsg{
			EventEum: enums.CertifyAlertContact,
			Data:     p.ActiveCode,
		}
		return tools.ToString(msg), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *AlertContactService) Select(param forms.AlertContactParam) *forms.AlertContactFormPage {
	param.TenantId = "1"
	db := global.DB
	return s.dao.Select(db, param)
}

func (s *AlertContactService) insertAlertContact(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContact, error) {
	//参数校验
	if p.ContactName == "" {
		return nil, errors.NewBusinessError("联系人名字不能为空")
	}
	//每个账号限制创建100个联系人
	var count int64
	global.DB.Model(&models.AlertContact{}).Where("tenant_id = ?", p.TenantId).Count(&count)
	if count >= constants.MaxContactNum {
		return nil, errors.NewBusinessError("联系人限制创建" + strconv.Itoa(constants.MaxContactNum) + "个")
	}
	//每个联系人最多加入5个联系组
	if len(p.GroupIdList) >= constants.MaxContactGroup {
		return nil, errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constants.MaxContactGroup) + "个联系组")
	}
	id := strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	p.ContactId = id
	//数据入库
	alertContact := &models.AlertContact{
		Id:          id,
		TenantId:    p.TenantId,
		Name:        p.ContactName,
		Description: p.Description,
		CreateUser:  p.CreateUser,
		Status:      1,
	}
	s.dao.Insert(db, alertContact)
	return alertContact, nil
}

func (s *AlertContactService) Update(param forms.AlertContactParam) error {
	db := global.DB
	if param.ContactName == "" {
		return errors.NewBusinessError("联系人名字不能为空")
	}
	//每个联系人最多加入5个联系组
	if len(param.GroupIdList) >= constants.MaxContactGroup {
		return errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constants.MaxContactGroup) + "个联系组")
	}
	currentTime := tools.GetNowStr()
	var alertContact = &models.AlertContact{
		Id:          param.ContactId,
		TenantId:    param.TenantId,
		Name:        param.ContactName,
		Status:      1,
		Description: param.Description,
		UpdateTime:  currentTime,
	}
	s.dao.Update(db, alertContact)
	return nil
}

func (s *AlertContactService) Delete(param forms.AlertContactParam) error {
	db := global.DB
	if param.ContactId == "" {
		return errors.NewBusinessError("联系人ID不能为空")
	}
	var alertContact = &models.AlertContact{
		Id:       param.ContactId,
		TenantId: param.TenantId,
	}
	s.dao.Delete(db, alertContact)
	return nil
}

func (s *AlertContactService) CertifyAlertContact(activeCode string) string {
	return s.dao.CertifyAlertContact(activeCode)
}
