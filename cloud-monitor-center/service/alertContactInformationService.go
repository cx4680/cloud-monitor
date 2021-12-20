package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
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
		List := s.updateAlertContactInformation(db, p)
		return tools.ToString(forms.MqMsg{
			EventEum: enums.UpdateAlertContactInformation,
			Data:     List,
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
	activeCode, isCertify := s.getActiveCode()
	var infoList []*models.AlertContactInformation
	if p.Phone != "" {
		alertContactInformationPhone := &models.AlertContactInformation{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   p.TenantId,
			ContactId:  p.ContactId,
			No:         p.Phone,
			Type:       1,
			IsCertify:  isCertify,
			ActiveCode: activeCode,
			CreateUser: p.CreateUser,
		}
		infoList = append(infoList, alertContactInformationPhone)
	}
	if p.Email != "" {
		alertContactInformationEmail := &models.AlertContactInformation{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   p.TenantId,
			ContactId:  p.ContactId,
			No:         p.Email,
			Type:       2,
			IsCertify:  isCertify,
			ActiveCode: activeCode,
			CreateUser: p.CreateUser,
		}
		infoList = append(infoList, alertContactInformationEmail)
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
	var messageSendDTO = messageCenter.MessageSendDTO{
		SourceType:   messageCenter.VERIFY,
		SenderId:     tenantId,
		Content:      tools.ToString(params),
		MsgEventCode: messageCenter.AddAlarmContactMail,
	}
	if noType == 1 {
		messageSendDTO.Target = []messageCenter.MessageTargetDTO{{
			Addr: no,
			Type: messageCenter.Phone,
		}}
	} else {
		messageSendDTO.Target = []messageCenter.MessageTargetDTO{{
			Addr: no,
			Type: messageCenter.Email,
		}}
	}
	err := messageCenter.NewService().Send(messageSendDTO)
	if err != nil {
		logger.Logger().Error(err)
	}
}
