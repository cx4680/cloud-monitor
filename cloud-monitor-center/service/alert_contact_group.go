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
	p := param.(form.AlertContactParam)
	//每个联系人最多加入5个联系组
	if len(p.GroupIdList) >= constant.MaxContactGroup {
		return "", errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MaxContactGroup) + "个联系组")
	}
	switch p.EventEum {
	case enum.InsertAlertContactGroup:
		//参数校验
		if p.GroupName == "" {
			return "", errors.NewBusinessError("联系组名字不能为空")
		}
		if len(p.ContactIdList) == 0 {
			return "", errors.NewBusinessError("请至少选择一位联系人")
		}
		//联系组限制创建10个
		var groupCount int64
		global.DB.Model(&model.AlertContactGroup{}).Where("tenant_id = ?", p.TenantId).Count(&groupCount)
		if groupCount >= constant.MaxGroupNum {
			return "", errors.NewBusinessError("联系组限制创建" + strconv.Itoa(constant.MaxGroupNum) + "个")
		}
		//联系组名不可重复
		var count int64
		global.DB.Model(&model.AlertContactGroup{}).Where("tenant_id = ? AND name = ?", p.TenantId, p.GroupName).Count(&count)
		if count >= 1 {
			return "", errors.NewBusinessError("联系组名重复")
		}
		alertContactGroup, err := s.insertAlertContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//保存联系人组关联
		p.GroupId = alertContactGroup.Id
		if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sys_rocketmq.AlertContactGroupTopic, p); err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.InsertAlertContactGroup,
			Data:     alertContactGroup,
		}), nil
	case enum.UpdateAlertContactGroup:
		//参数校验
		if p.GroupName == "" {
			return "", errors.NewBusinessError("联系组名字不能为空")
		}
		//联系组名不可重复
		var count int64
		global.DB.Model(&model.AlertContactGroup{}).Where("tenant_id = ? AND id != ? AND name = ?", p.TenantId, p.GroupId, p.GroupName).Count(&count)
		if count >= 1 {
			return "", errors.NewBusinessError("联系组名重复")
		}
		alertContactGroup, err := s.updateAlertContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sys_rocketmq.AlertContactGroupTopic, p); err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.UpdateAlertContactGroup,
			Data:     alertContactGroup,
		}), nil
	case enum.DeleteAlertContactGroup:
		alertContactGroup, err := s.deleteAlertContactGroup(db, p)
		if err != nil {
			return "", err
		}
		//更新联系人组关联
		if err := s.alertContactGroupRelService.PersistenceInner(db, s.alertContactGroupRelService, sys_rocketmq.AlertContactGroupTopic, p); err != nil {
			return "", err
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.DeleteAlertContactGroup,
			Data:     alertContactGroup,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")
	}
}

func (s *AlertContactGroupService) SelectAlertContactGroup(param form.AlertContactParam) *form.AlertContactFormPage {
	return s.dao.SelectAlertContactGroup(global.DB, param)
}

func (s *AlertContactGroupService) SelectAlertGroupContact(param form.AlertContactParam) *form.AlertContactFormPage {
	return s.dao.SelectAlertGroupContact(global.DB, param)
}

func (s *AlertContactGroupService) insertAlertContactGroup(db *gorm.DB, p form.AlertContactParam) (*model.AlertContactGroup, error) {
	currentTime := util.GetNow()
	alertContactGroup := &model.AlertContactGroup{
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

func (s *AlertContactGroupService) updateAlertContactGroup(db *gorm.DB, p form.AlertContactParam) (*model.AlertContactGroup, error) {
	var oldAlertContactGroup = &model.AlertContactGroup{}
	global.DB.Where("tenant_id = ? AND id = ?", p.TenantId, p.GroupId).First(oldAlertContactGroup)
	currentTime := util.GetNow()
	var alertContactGroup = &model.AlertContactGroup{
		Id:          p.GroupId,
		TenantId:    p.TenantId,
		Name:        p.GroupName,
		Description: p.Description,
		UpdateTime:  currentTime,
		CreateTime:  oldAlertContactGroup.CreateTime,
		CreateUser:  oldAlertContactGroup.CreateUser,
	}
	s.dao.Update(db, alertContactGroup)
	return alertContactGroup, nil
}

func (s *AlertContactGroupService) deleteAlertContactGroup(db *gorm.DB, p form.AlertContactParam) (*model.AlertContactGroup, error) {
	var alertContactGroup = &model.AlertContactGroup{
		Id:       p.GroupId,
		TenantId: p.TenantId,
	}
	s.dao.Delete(db, alertContactGroup)
	return alertContactGroup, nil
}
