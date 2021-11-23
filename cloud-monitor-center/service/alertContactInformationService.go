package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	activeCode, isCertify := s.getActiveCode()
	alertContactInformationPhone := &models.AlertContactInformation{
		TenantId:   p.TenantId,
		ContactId:  p.ContactId,
		No:         p.Phone,
		Type:       1,
		IsCertify:  isCertify,
		ActiveCode: activeCode,
		CreateUser: p.CreateUser,
	}

	alertContactInformationEmail := &models.AlertContactInformation{
		TenantId:   p.TenantId,
		ContactId:  p.ContactId,
		No:         p.Email,
		Type:       2,
		IsCertify:  isCertify,
		ActiveCode: activeCode,
		CreateUser: p.CreateUser,
	}
	infoList := []*models.AlertContactInformation{alertContactInformationPhone, alertContactInformationEmail}

	s.dao.InsertBatch(db, infoList)

	return tools.ToString(forms.MqMsg{
		EventEum: enums.InsertAlertContactInformation,
		Data:     infoList,
	}), nil
}

func (s *AlertContactInformationService) getActiveCode() (string, int) {
	if !config.GetCommonConfig().HasNoticeModel {
		return "", 1
	}
	return strings.ReplaceAll(uuid.New().String(), "-", ""), 0
}
