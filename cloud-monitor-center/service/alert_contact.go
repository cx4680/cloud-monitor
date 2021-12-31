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
	p := param.(form.AlertContactParam)
	//每个联系人最多加入5个联系组
	if len(p.GroupIdList) > constant.MaxContactGroup {
		return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
	}
	switch p.EventEum {
	case enum.InsertAlertContact:
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
		msg := form.MqMsg{
			EventEum: enum.InsertAlertContact,
			Data:     alertContact,
		}
		return jsonutil.ToString(msg), nil
	case enum.UpdateAlertContact:
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
		msg := form.MqMsg{
			EventEum: enum.UpdateAlertContact,
			Data:     alertContact,
		}
		return jsonutil.ToString(msg), nil
	case enum.DeleteAlertContact:
		alertContact, err := s.deleteAlertContact(db, p)
		if err != nil {
			return "", err
		}
		//删除 联系方式 和 联系人组关联
		if err := s.persistenceInner(db, p); err != nil {
			return "", err
		}
		msg := form.MqMsg{
			EventEum: enum.DeleteAlertContact,
			Data:     alertContact,
		}
		return jsonutil.ToString(msg), nil
	case enum.CertifyAlertContact:
		s.dao.CertifyAlertContact(db, p.ActiveCode)
		msg := form.MqMsg{
			EventEum: enum.CertifyAlertContact,
			Data:     p.ActiveCode,
		}
		return jsonutil.ToString(msg), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *AlertContactService) SelectAlertContact(param form.AlertContactParam) *form.AlertContactFormPage {
	return s.dao.Select(global.DB, param)
}

func (s *AlertContactService) insertAlertContact(db *gorm.DB, p form.AlertContactParam) (*model.AlertContact, error) {
	//每个账号限制创建100个联系人
	var count int64
	global.DB.Model(&model.AlertContact{}).Where("tenant_id = ?", p.TenantId).Count(&count)
	if count >= constant.MaxContactNum {
		return nil, errors.NewBusinessError("联系人限制创建" + strconv.Itoa(constant.MaxContactNum) + "个")
	}
	currentTime := util.GetNowStr()
	alertContact := &model.AlertContact{
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

func (s *AlertContactService) updateAlertContact(db *gorm.DB, p form.AlertContactParam) (*model.AlertContact, error) {
	//校验联系人是否为该租户所有
	var check int64
	db.Model(&model.AlertContact{}).Where("tenant_id = ? AND id = ?", p.TenantId, p.ContactId).Count(&check)
	if check == 0 {
		return nil, errors.NewBusinessError("该租户无此联系人")
	}

	currentTime := util.GetNowStr()
	var alertContact = &model.AlertContact{
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

func (s *AlertContactService) deleteAlertContact(db *gorm.DB, p form.AlertContactParam) (*model.AlertContact, error) {
	//校验联系人是否为该租户所有
	var check int64
	db.Model(&model.AlertContact{}).Where("tenant_id = ? AND contact_id = ?", p.TenantId, p.ContactId).Count(&check)
	if p.ContactId == "" {
		return nil, errors.NewBusinessError("联系人ID不能为空")
	}
	var alertContact = &model.AlertContact{
		Id:       p.ContactId,
		TenantId: p.TenantId,
	}
	s.dao.Delete(db, alertContact)
	return alertContact, nil
}

func (s *AlertContactService) persistenceInner(db *gorm.DB, p form.AlertContactParam) error {
	//联系人组关联
	if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sys_rocketmq.AlertContactGroupTopic, p); err != nil {
		return err
	}
	//联系方式
	if err := s.alertContactInformationService.PersistenceInner(db, s.alertContactInformationService, sys_rocketmq.AlertContactTopic, p); err != nil {
		return err
	}
	return nil
}

func (s *AlertContactService) GetTenantId(activeCode string) string {
	var model = &model.AlertContactInformation{}
	global.DB.Where("active_code = ?", activeCode).First(model)
	return model.TenantId
}
