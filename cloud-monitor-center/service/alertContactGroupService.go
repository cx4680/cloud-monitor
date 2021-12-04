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

type AlertContactGroupService struct {
	service.AbstractSyncServiceImpl
	dao                         *dao.AlertContactGroupDao
	alertContactGroupRelService *AlertContactGroupRelService
}

func NewAlertContactGroupService(alertContactGroupRelService *AlertContactGroupRelService) *AlertContactGroupService {
	return &AlertContactGroupService{
		AbstractSyncServiceImpl:     service.AbstractSyncServiceImpl{},
		dao:                         dao.AlertContactGroup,
		alertContactGroupRelService: alertContactGroupRelService,
	}
}

func (s *AlertContactGroupService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(forms.AlertContactParam)
	//每个联系人最多加入5个联系组
	if len(p.GroupIdList) >= constants.MaxContactGroup {
		return "", errors.NewBusinessError("每个联系组最多加入" + strconv.Itoa(constants.MaxGroupNum) + "个联系组")
	}
	switch p.EventEum {
	case enums.InsertAlertContactGroup:
		//参数校验
		if p.GroupName == "" {
			return "", errors.NewBusinessError("联系组名字不能为空")
		}
		//联系组名不可重复
		var count int64
		global.DB.Model(&models.AlertContactGroup{}).Where("tenant_id = ? AND name = ?", p.TenantId, p.GroupName).Count(&count)
		if count >= 1 {
			return "", errors.NewBusinessError("联系组名重复")
		}
		alertContactGroup, err := s.insertAlertContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//保存联系人组关联
		p.GroupId = alertContactGroup.Id
		if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sysRocketMq.AlertContactGroupTopic, p); err != nil {
			return "", err
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.InsertAlertContactGroup,
			Data:     alertContactGroup,
		}), nil
	case enums.UpdateAlertContactGroup:
		//参数校验
		if p.GroupName == "" {
			return "", errors.NewBusinessError("联系组名字不能为空")
		}
		//联系组名不可重复
		var count int64
		global.DB.Model(&models.AlertContactGroup{}).Where("tenant_id = ? AND id != ? AND name = ?", p.TenantId, p.GroupId, p.GroupName).Count(&count)
		if count >= 1 {
			return "", errors.NewBusinessError("联系组名重复")
		}
		alertContactGroup, err := s.updateAlertContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sysRocketMq.AlertContactGroupTopic, p); err != nil {
			return "", err
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.UpdateAlertContactGroup,
			Data:     alertContactGroup,
		}), nil
	case enums.DeleteAlertContactGroup:
		alertContactGroup, err := s.deleteAlertContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sysRocketMq.AlertContactGroupTopic, p); err != nil {
			return "", err
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.DeleteAlertContactGroup,
			Data:     alertContactGroup,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *AlertContactGroupService) SelectAlertContactGroup(param forms.AlertContactParam) *[]forms.AlertContactGroupForm {
	param.TenantId = "1"
	db := global.DB
	return s.dao.SelectAlertContactGroup(db, param)
}

func (s *AlertContactGroupService) SelectAlertGroupContact(param forms.AlertContactParam) *[]forms.AlertContactForm {
	param.TenantId = "1"
	db := global.DB
	return s.dao.SelectAlertGroupContact(db, param)
}

func (s *AlertContactGroupService) insertAlertContactGroup(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContactGroup, error) {
	currentTime := tools.GetNowStr()
	alertContactGroup := &models.AlertContactGroup{
		Id:          strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		TenantId:    p.TenantId,
		Name:        p.GroupName,
		Description: p.Description,
		CreateUser:  p.CreateUser,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
	}
	s.dao.Insert(db, alertContactGroup)
	return alertContactGroup, nil
}

func (s *AlertContactGroupService) updateAlertContactGroup(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContactGroup, error) {
	currentTime := tools.GetNowStr()
	var alertContactGroup = &models.AlertContactGroup{
		Id:          p.GroupId,
		TenantId:    p.TenantId,
		Name:        p.GroupName,
		Description: p.Description,
		UpdateTime:  currentTime,
	}
	s.dao.Update(db, alertContactGroup)
	return alertContactGroup, nil
}

func (s *AlertContactGroupService) deleteAlertContactGroup(db *gorm.DB, p forms.AlertContactParam) (*models.AlertContactGroup, error) {
	var alertContactGroup = &models.AlertContactGroup{
		Id:       p.GroupId,
		TenantId: p.TenantId,
	}
	s.dao.Delete(db, alertContactGroup)
	return alertContactGroup, nil
}
