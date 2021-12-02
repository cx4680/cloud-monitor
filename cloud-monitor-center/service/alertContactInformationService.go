package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
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
		infoList := s.getInformationList(db, p)
		return tools.ToString(forms.MqMsg{
			EventEum: enums.InsertAlertContactInformation,
			Data:     infoList,
		}), nil
	case enums.UpdateAlertContact:
		p.EventEum = enums.DeleteAlertContactInformation
		s.PersistenceLocal(db, p)
		p.EventEum = enums.InsertAlertContactInformation
		s.PersistenceLocal(db, p)
		return "", nil
		//infoList := s.getInformationList(p)
		//var alertContactInformation = &models.AlertContactInformation{
		//	Id:       p.ContactId,
		//	TenantId: p.TenantId,
		//}
		//s.dao.Update(db, infoList, alertContactInformation)
		//return tools.ToString(forms.MqMsg{
		//	EventEum: enums.UpdateAlertContactGroup,
		//	Data:     alertContactInformation,
		//}), nil
	case enums.DeleteAlertContact:
		var alertContactInformation = &models.AlertContactInformation{
			ContactId: p.ContactId,
			TenantId:  p.TenantId,
		}
		s.dao.Delete(db, alertContactInformation)
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

func (s *AlertContactInformationService) getInformationList(db *gorm.DB, p forms.AlertContactParam) []*models.AlertContactInformation {
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
	s.dao.InsertBatch(db, infoList)
	return infoList
}
