package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type AlertContactInformationService struct {
	service.AbstractSyncServiceImpl
	messageSvc *service.MessageService
	dao        *dao.AlertContactInformationDao
}

func NewAlertContactInformationService(messageSvc *service.MessageService) *AlertContactInformationService {
	return &AlertContactInformationService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		messageSvc:              messageSvc,
		dao:                     dao.AlertContactInformation,
	}
}

// PersistenceLocal 插入两条数据
func (s *AlertContactInformationService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(form.AlertContactParam)
	switch p.EventEum {
	case enum.InsertAlertContact:
		list := s.insertAlertContactInformation(db, p)
		for _, information := range list {
			s.sendCertifyMsg(information.TenantId, information.No, information.Type, information.ActiveCode, information.ContactId)
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.InsertAlertContactInformation,
			Data:     list,
		}), nil
	case enum.UpdateAlertContact:
		list := s.updateAlertContactInformation(db, p)
		for _, information := range list {
			s.sendCertifyMsg(information.TenantId, information.No, information.Type, information.ActiveCode, information.ContactId)
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.UpdateAlertContactInformation,
			Data:     list,
		}), nil
	case enum.DeleteAlertContact:
		alertContactInformation := s.deleteAlertContactInformation(db, p)
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.DeleteAlertContactInformation,
			Data:     alertContactInformation,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")

	}

}

func (s *AlertContactInformationService) getActiveCode(noType uint8) (string, uint8) {
	if config.Cfg.Common.MsgIsOpen == config.MsgClose {
		return "", 1
	}
	for k := range global.NoticeChannelMap {
		if (k == config.MsgChannelSms && noType == constant.Phone) || (k == config.MsgChannelEmail && noType == constant.Email) {
			return strings.ReplaceAll(uuid.New().String(), "-", ""), 0
		}
	}
	return "", 1
}

func (s *AlertContactInformationService) insertAlertContactInformation(db *gorm.DB, p form.AlertContactParam) []*model.AlertContactInformation {
	List := s.getInformationList(p)
	s.dao.InsertBatch(db, List)
	return List
}

func (s *AlertContactInformationService) updateAlertContactInformation(db *gorm.DB, p form.AlertContactParam) []*model.AlertContactInformation {
	List := s.getInformationList(p)
	s.dao.Update(db, List)
	return List
}

func (s *AlertContactInformationService) deleteAlertContactInformation(db *gorm.DB, p form.AlertContactParam) *model.AlertContactInformation {
	var alertContactInformation = &model.AlertContactInformation{
		ContactId: p.ContactId,
		TenantId:  p.TenantId,
	}
	s.dao.Delete(db, alertContactInformation)
	return alertContactInformation
}

func (s *AlertContactInformationService) getInformationList(p form.AlertContactParam) []*model.AlertContactInformation {
	var infoList []*model.AlertContactInformation
	if strutil.IsNotBlank(p.Phone) {
		if information := s.buildInformation(p, p.Phone, 1); information != nil {
			infoList = append(infoList, information)
		}
	}
	if strutil.IsNotBlank(p.Email) {
		if information := s.buildInformation(p, p.Email, 2); information != nil {
			infoList = append(infoList, information)
		}
	}
	return infoList
}

//发送激活信息
func (s *AlertContactInformationService) sendCertifyMsg(tenantId string, no string, noType uint8, activeCode string, contactId string) {
	if no == "" {
		return
	}
	params := make(map[string]string)
	params["userName"] = service.NewTenantService().GetTenantInfo(tenantId).Name
	params["verifyBtn"] = config.Cfg.Common.CertifyInformationUrl + activeCode
	params["activationdomain"] = config.Cfg.Common.CertifyInformationUrl + activeCode

	var t message_center.ReceiveType
	if noType == constant.Phone {
		t = message_center.Phone
	} else if noType == constant.Email {
		t = message_center.Email
	}
	s.messageSvc.SendCertifyMsg(message_center.MessageSendDTO{
		SenderId:   tenantId,
		Type:       t,
		SourceType: message_center.VERIFY,
		Targets:    []string{no},
		Content:    jsonutil.ToString(params),
	}, contactId)
}

func (s *AlertContactInformationService) buildInformation(p form.AlertContactParam, no string, noType uint8) *model.AlertContactInformation {
	activeCode, isCertify := s.getActiveCode(noType)
	var information = &model.AlertContactInformation{}
	//判断新增的联系方式是否已存在，若存在则不修改，若不存在，则删除旧号码，添加新号码
	var count int64
	global.DB.Model(&model.AlertContactInformation{}).Where("tenant_id = ? AND contact_id = ? AND no = ? AND type = ?", p.TenantId, p.ContactId, no, noType).Count(&count)
	if count == 0 {
		information = &model.AlertContactInformation{
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
