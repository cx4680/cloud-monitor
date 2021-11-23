package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
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
	//数据入库
	alertContact := &models.AlertContact{
		TenantId:    p.TenantId,
		Name:        p.ContactName,
		Description: p.Description,
		CreateUser:  p.CreateUser,
	}
	s.dao.Insert(db, alertContact)

	//保存联系人组
	if err := s.alertContactInformationService.PersistenceInner(db, s.alertContactInformationService, config.GetRocketmqConfig().AlertContactTopic, param); err != nil {
		return "", errors.NewBusinessError(err.Error())
	}
	//保存联系人组关系
	if err := s.alertContactGroupService.PersistenceInner(db, s.alertContactGroupService, config.GetRocketmqConfig().AlertContactTopic, param); err != nil {
		return "", errors.NewBusinessError(err.Error())
	}

	msg := forms.MqMsg{
		EventEum: enums.InsertAlertContact,
		Data:     alertContact,
	}
	return tools.ToString(msg), nil
}
