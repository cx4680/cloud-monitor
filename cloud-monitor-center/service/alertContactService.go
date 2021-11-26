package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
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
}

func NewAlertContactService(alertContactGroupService *AlertContactGroupService, alertContactInformationService *AlertContactInformationService) *AlertContactService {
	return &AlertContactService{
		AbstractSyncServiceImpl:        service.AbstractSyncServiceImpl{},
		dao:                            dao.AlertContact,
		alertContactGroupService:       alertContactGroupService,
		alertContactInformationService: alertContactInformationService,
	}
}

func (s *AlertContactService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(forms.AlertContactParam)
	switch p.EventEum {
	case enums.InsertAlertContact:
		//参数校验
		if p.ContactName == "" {
			return "", errors.NewBusinessError("联系人名字不能为空")
		}
		//每个账号限制创建100个联系人
		var count int64
		global.DB.Model(&models.AlertContact{}).Where("tenant_id = ?", p.TenantId).Count(&count)
		if count >= constant.MAX_CONTACT_NUM {
			return "", errors.NewBusinessError("联系人限制创建" + strconv.Itoa(constant.MAX_CONTACT_NUM) + "个")
		}
		//每个联系人最多加入5个联系组
		if len(p.GroupIdList) >= constant.MAX_CONTACT_GROUP {
			return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
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

		//保存联系方式
		if err := s.alertContactInformationService.PersistenceInner(db, s.alertContactInformationService, config.GetRocketmqConfig().AlertContactTopic, param); err != nil {
			return "", errors.NewBusinessError(err.Error())
		}
		//保存联系人组关联
		if err := s.alertContactGroupService.PersistenceInner(db, s.alertContactGroupService, config.GetRocketmqConfig().AlertContactTopic, param); err != nil {
			return "", errors.NewBusinessError(err.Error())
		}
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
		if len(p.GroupIdList) >= constant.MAX_CONTACT_GROUP {
			return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
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
		if err := s.alertContactInformationService.PersistenceInner(db, s.alertContactInformationService, config.GetRocketmqConfig().AlertContactTopic, param); err != nil {
			return "", errors.NewBusinessError(err.Error())
		}
		//更新联系人组关联
		if err := s.alertContactGroupService.PersistenceInner(db, s.alertContactGroupService, config.GetRocketmqConfig().AlertContactTopic, param); err != nil {
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

func (s *AlertContactService) Insert(param forms.AlertContactParam) error {
	db := global.DB
	if param.ContactName == "" {
		return errors.NewBusinessError("联系人名字不能为空")
	}
	//每个账号限制创建100个联系人
	var count int64
	global.DB.Model(&models.AlertContact{}).Where("tenant_id = ?", param.TenantId).Count(&count)
	if count >= constant.MAX_CONTACT_NUM {
		return errors.NewBusinessError("联系人限制创建" + strconv.Itoa(constant.MAX_CONTACT_NUM) + "个")
	}
	//每个联系人最多加入5个联系组
	if len(param.GroupIdList) >= constant.MAX_CONTACT_GROUP {
		return errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
	}

	currentTime := tools.GetNowStr()
	contactId := strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	param.ContactId = contactId
	var alertContact = &models.AlertContact{
		Id:          contactId,
		TenantId:    param.TenantId,
		Name:        param.ContactName,
		Status:      1,
		Description: param.Description,
		CreateUser:  param.CreateUser,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
	}
	s.dao.Insert(db, alertContact)
	return nil
}

func (s *AlertContactService) Update(param forms.AlertContactParam) error {
	db := global.DB
	if param.ContactName == "" {
		return errors.NewBusinessError("联系人名字不能为空")
	}
	//每个联系人最多加入5个联系组
	if len(param.GroupIdList) >= constant.MAX_CONTACT_GROUP {
		return errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
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
