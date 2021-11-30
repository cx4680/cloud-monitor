package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactGroupService struct {
	service.AbstractSyncServiceImpl
	dao             *dao.AlertContactGroupDao
	alertContactDao *dao.AlertContactDao
}

func NewAlertContactGroupService() *AlertContactGroupService {
	return &AlertContactGroupService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		dao:                     dao.AlertContactGroup,
		alertContactDao:         dao.AlertContact,
	}
}

func (acgs *AlertContactGroupService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(forms.AlertContactParam)
	switch p.EventEum {
	case enums.InsertAlertContact:

	}
	relList := make([]*models.AlertContactGroupRel, len(p.GroupIdList))
	var count int64
	for i, v := range p.GroupIdList {
		db.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ?", p.TenantId).Where("group_id = ?", v).Count(&count)
		if count >= constants.MaxContactNum {
			return "", errors.NewBusinessError("每组联系人限制创建" + strconv.Itoa(constants.MaxContactNum) + "个")
		}
		var alertContactGroupRel = &models.AlertContactGroupRel{
			TenantId:   p.TenantId,
			ContactId:  p.ContactId,
			GroupId:    v,
			CreateUser: p.CreateUser,
		}
		relList[i] = alertContactGroupRel
	}

	acgs.dao.InsertGroupRelBatch(db, relList)
	return tools.ToString(relList), nil
}

func insert() {

}
