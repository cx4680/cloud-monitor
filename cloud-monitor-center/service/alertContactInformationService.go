package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type AlertContactInformationService struct {
	service.AbstractSyncServiceImpl
	dao *dao.AlertContactInformationDao
}

func NewAlertContactInformationService() *AlertContactInformationService {
	return &AlertContactInformationService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		dao:                     dao.AlertContactInformation,
	}
}

// PersistenceLocal 插入两条数据
func (s *AlertContactInformationService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(forms.AlertContactParam)
	switch p.EventEum {
	case enums.InsertAlertContact:
		list := s.insertAlertContactInformation(db, p)
		for _, information := range list {
			sendCertifyMsg(information.TenantId, information.ContactId, information.No, information.Type, information.ActiveCode)
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.InsertAlertContactInformation,
			Data:     list,
		}), nil
	case enums.UpdateAlertContact:
		list := s.updateAlertContactInformation(db, p)
		for _, information := range list {
			sendCertifyMsg(information.TenantId, information.ContactId, information.No, information.Type, information.ActiveCode)
		}
		return tools.ToString(forms.MqMsg{
			EventEum: enums.UpdateAlertContactInformation,
			Data:     list,
		}), nil
	case enums.DeleteAlertContact:
		alertContactInformation := s.deleteAlertContactInformation(db, p)
		return tools.ToString(forms.MqMsg{
			EventEum: enums.DeleteAlertContactInformation,
			Data:     alertContactInformation,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")

	}

}

func (s *AlertContactInformationService) getActiveCode() (string, int) {
	if !config.GetCommonConfig().HasNoticeModel {
		return "", 1
	}
	return strings.ReplaceAll(uuid.New().String(), "-", ""), 0
}

func (s *AlertContactInformationService) insertAlertContactInformation(db *gorm.DB, p forms.AlertContactParam) []*models.AlertContactInformation {
	List := s.getInformationList(p)
	s.dao.InsertBatch(db, List)
	return List
}

func (s *AlertContactInformationService) updateAlertContactInformation(db *gorm.DB, p forms.AlertContactParam) []*models.AlertContactInformation {
	List := s.getInformationList(p)
	s.dao.Update(db, List)
	return List
}

func (s *AlertContactInformationService) deleteAlertContactInformation(db *gorm.DB, p forms.AlertContactParam) *models.AlertContactInformation {
	var alertContactInformation = &models.AlertContactInformation{
		ContactId: p.ContactId,
		TenantId:  p.TenantId,
	}
	s.dao.Delete(db, alertContactInformation)
	return alertContactInformation
}

func (s *AlertContactInformationService) getInformationList(p forms.AlertContactParam) []*models.AlertContactInformation {
	var infoList []*models.AlertContactInformation
	if tools.IsNotBlank(p.Phone) {
		if information := s.buildInformation(p, p.Phone, 1); information != nil {
			infoList = append(infoList, information)
		}
	}
	if tools.IsNotBlank(p.Email) {
		if information := s.buildInformation(p, p.Email, 2); information != nil {
			infoList = append(infoList, information)
		}
	}
	return infoList
}

//发送激活信息
func sendCertifyMsg(tenantId string, contactId string, no string, noType int, activeCode string) {
	if no == "" {
		return
	}
	params := make(map[string]string)
	params["userName"] = service.NewTenantService().GetTenantInfo(tenantId).Name
	params["verifyBtn"] = config.GetCommonConfig().CertifyInformationUrl + activeCode
	params["activationdomain"] = config.GetCommonConfig().CertifyInformationUrl + activeCode

	var t messageCenter.ReceiveType
	if noType == 1 {
		t = messageCenter.Phone
	} else {
		t = messageCenter.Email
	}

	messageSendDTO := messageCenter.MessageSendDTO{
		SenderId:   tenantId,
		Type:       t,
		SourceType: messageCenter.VERIFY,
		Targets:    []string{no},
		Content:    tools.ToString(params),
	}
	err := messageCenter.NewService().Send(messageSendDTO)
	if err != nil {
		logger.Logger().Error(err)
	}
}

func (s *AlertContactInformationService) buildInformation(p forms.AlertContactParam, no string, noType int) *models.AlertContactInformation {
	activeCode, isCertify := s.getActiveCode()
	var information = &models.AlertContactInformation{}
	//判断新增的联系方式是否已存在，若存在则不修改，若不存在，则删除旧号码，添加新号码
	var count int64
	global.DB.Model(&models.AlertContactInformation{}).Where("tenant_id = ? AND contact_id = ? AND no = ? AND type = ?", p.TenantId, p.ContactId, no, noType).Count(&count)
	if count == 0 {
		information = &models.AlertContactInformation{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   p.TenantId,
			ContactId:  p.ContactId,
			No:         no,
			Type:       noType,
			IsCertify:  isCertify,
			ActiveCode: activeCode,
			CreateUser: p.CreateUser,
		}
		return information
	}
	return nil
}
