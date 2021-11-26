package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	forms2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	models2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactGroupService struct {
	service.AbstractSyncServiceImpl
	dao             *dao2.AlertContactGroupDao
	alertContactDao *dao2.AlertContactDao
}

func NewAlertContactGroupService() *AlertContactGroupService {
	return &AlertContactGroupService{
		AbstractSyncServiceImpl: service.AbstractSyncServiceImpl{},
		dao:                     dao2.AlertContactGroup,
		alertContactDao:         dao2.AlertContact,
	}
}

func (acgs *AlertContactGroupService) PersistenceLocal(db *gorm.DB, param interface{}) (string, error) {
	p := param.(forms2.AlertContactParam)

	relList := make([]*models2.AlertContactGroupRel, len(p.GroupIdList))
	var count int64
	for i, v := range p.GroupIdList {
		db.Model(&models2.AlertContactGroupRel{}).Where("tenant_id = ?", p.TenantId).Where("group_id = ?", v).Count(&count)
		if count >= constants.MAX_CONTACT_NUM {
			return "", errors.NewBusinessError("每组联系人限制创建" + strconv.Itoa(constants.MAX_CONTACT_NUM) + "个")
		}
		var alertContactGroupRel = &models2.AlertContactGroupRel{
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
