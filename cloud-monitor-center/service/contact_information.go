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
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type ContactInformationService struct {
	service.AbstractSyncServiceImpl
	messageSvc *service.MessageService
	dao        *dao.ContactInformationDao
}

func NewContactInformationService(messageSvc *service.MessageService) *ContactInformationService {
	return &ContactInformationService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		messageSvc:              messageSvc,
		dao:                     dao.ContactInformation,
	}
}

// PersistenceLocal 插入两条数据
func (s *ContactInformationService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(form.ContactParam)
	switch p.EventEum {
	case enum.InsertContact:
		err := checkAddressSize(p.Phone, p.Email)
		if err != nil {
			return "", err
		}
		list := s.insertContactInformation(db, p)
		for _, information := range list {
			if strutil.IsNotBlank(information.ActiveCode) {
				s.sendActivateMsg(information.TenantId, information.Address, information.Type, information.ActiveCode, information.ContactBizId)
			}
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.InsertContactInformation,
			Data:     list,
		}), nil
	case enum.UpdateContact:
		err := checkAddressSize(p.Phone, p.Email)
		if err != nil {
			return "", err
		}
		list := s.updateContactInformation(db, p)
		for _, information := range list {
			s.sendActivateMsg(information.TenantId, information.Address, information.Type, information.ActiveCode, information.ContactBizId)
		}
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.UpdateContactInformation,
			Data:     list,
		}), nil
	case enum.DeleteContact:
		contactInformation := s.deleteContactInformation(db, p)
		return jsonutil.ToString(form.MqMsg{
			EventEum: enum.DeleteContactInformation,
			Data:     contactInformation,
		}), nil
	default:
		return "", errors.NewBusinessError("系统异常")

	}

}

func (s *ContactInformationService) getActiveCode(addressType uint8) (string, uint8) {
	if config.Cfg.Common.EnvType == config.ProprietaryCloud || config.Cfg.Common.MsgIsOpen == config.MsgClose {
		return "", constant.Activated
	}
	for _, v := range global.NoticeChannelList {
		if (v.Code == config.MsgChannelSms && addressType == constant.Phone) || (v.Code == config.MsgChannelEmail && addressType == constant.Email) {
			return strings.ReplaceAll(uuid.New().String(), "-", ""), constant.NotActive
		}
	}
	return "", constant.Activated
}

func (s *ContactInformationService) insertContactInformation(db *gorm.DB, p form.ContactParam) []*model.ContactInformation {
	list := s.getInformationList(p)
	s.dao.InsertBatch(db, list)
	return list
}

func (s *ContactInformationService) updateContactInformation(db *gorm.DB, p form.ContactParam) []*model.ContactInformation {
	list := s.getInformationList(p)
	s.dao.Update(db, list)
	return list
}

func (s *ContactInformationService) deleteContactInformation(db *gorm.DB, p form.ContactParam) *model.ContactInformation {
	var contactInformation = &model.ContactInformation{
		ContactBizId: p.ContactBizId,
		TenantId:     p.TenantId,
	}
	s.dao.Delete(db, contactInformation)
	return contactInformation
}

func (s *ContactInformationService) getInformationList(p form.ContactParam) []*model.ContactInformation {
	var infoList []*model.ContactInformation
	if information := s.buildInformation(p, p.Phone, constant.Phone); information != nil {
		infoList = append(infoList, information)
	}
	if information := s.buildInformation(p, p.Email, constant.Email); information != nil {
		infoList = append(infoList, information)
	}
	return infoList
}

//发送激活信息
func (s *ContactInformationService) sendActivateMsg(tenantId string, address string, addressType uint8, activeCode string, contactBizId string) {
	if strutil.IsBlank(address) {
		return
	}
	params := make(map[string]string)
	params["userName"] = service.NewTenantService().GetTenantInfo(tenantId).Name
	params["verifyBtn"] = config.Cfg.Common.ActivateInformationUrl + activeCode
	params["activationdomain"] = config.Cfg.Common.ActivateInformationUrl + activeCode

	var t message_center.ReceiveType
	if addressType == constant.Phone {
		t = message_center.Phone
	} else if addressType == constant.Email {
		t = message_center.Email
	}
	s.messageSvc.SendActivateMsg(message_center.MessageSendDTO{
		SenderId:   tenantId,
		Type:       t,
		SourceType: message_center.VERIFY,
		Targets:    []string{address},
		Content:    jsonutil.ToString(params),
	}, contactBizId)
}

func (s *ContactInformationService) buildInformation(p form.ContactParam, address string, addressType uint8) *model.ContactInformation {
	activeCode, state := s.getActiveCode(addressType)
	var information = &model.ContactInformation{}
	//判断新增的联系方式是否已存在，若存在则不修改，若不存在，则删除旧号码，添加新号码
	var count int64
	global.DB.Model(&model.ContactInformation{}).Where("tenant_id = ? AND contact_biz_id = ? AND address = ? AND type = ?", p.TenantId, p.ContactBizId, address, addressType).Count(&count)
	if count == 0 {
		information = &model.ContactInformation{
			TenantId:     p.TenantId,
			ContactBizId: p.ContactBizId,
			Address:      address,
			Type:         addressType,
			State:        state,
			ActiveCode:   activeCode,
			CreateUser:   p.CreateUser,
		}
		return information
	}
	return nil
}

func checkAddressSize(phone, email string) error {
	if strutil.IsNotBlank(phone) && len(phone) != constant.PhoneSize {
		return errors.NewBusinessError("手机号必须为" + strconv.Itoa(constant.PhoneSize) + "位")
	}
	if strutil.IsNotBlank(email) && len(email) > constant.MaxEmailSize {
		return errors.NewBusinessError("邮箱限制" + strconv.Itoa(constant.MaxEmailSize) + "个字符")
	}
	return nil
}
