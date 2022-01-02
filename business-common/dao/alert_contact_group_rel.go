package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"gorm.io/gorm"
)

type AlertContactGroupRelDao struct {
}

const (
	checkAlertContact = "SELECT " +
		"ac.id AS contact_id, " +
		"ac.name AS contact_name, " +
		"ANY_VALUE( acg.group_id ) AS group_id, " +
		"ANY_VALUE( acg.group_name ) AS group_name, " +
		"ANY_VALUE( acg.group_count ) AS group_count " +
		"FROM " +
		"alert_contact AS ac " +
		"LEFT JOIN alert_contact_information AS aci ON ac.id = aci.contact_id AND ac.tenant_id = aci.tenant_id " +
		"LEFT JOIN " +
		"( " +
		"SELECT " +
		"acgr.contact_id AS contact_id, " +
		"ANY_VALUE( acgr.tenant_id ) AS tenant_id, " +
		"GROUP_CONCAT( acg.id ) AS group_id, " +
		"GROUP_CONCAT( acg.name ) AS group_name, " +
		"COUNT( acgr.contact_id ) AS group_count  " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id AND acg.tenant_id = acgr.tenant_id " +
		"WHERE acg.id != ? " +
		"GROUP BY " +
		"acgr.contact_id " +
		") " +
		"AS acg ON ac.id = acg.contact_id AND ac.tenant_id = acg.tenant_id " +
		"WHERE " +
		"ac.status = 1 " +
		"AND ac.tenant_id = ? " +
		"AND ac.id = ? " +
		"GROUP BY " +
		"ac.id "

	checkAlertContactGroup = "SELECT " +
		"acg.id AS group_id, " +
		"acg.name AS group_name, " +
		"COUNT( acgr.group_id ) AS contact_count " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id AND acg.tenant_id = acgr.tenant_id " +
		"WHERE " +
		"acg.tenant_id = ? " +
		"AND acg.id = ? " +
		"GROUP BY " +
		"acg.id "
)

var AlertContactGroupRel = new(AlertContactGroupRelDao)

func (d *AlertContactGroupRelDao) Insert(db *gorm.DB, entity *model.AlertContactGroupRel) {
	currentTime := util.GetNow()
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (d *AlertContactGroupRelDao) InsertBatch(db *gorm.DB, list []*model.AlertContactGroupRel) {
	if len(list) == 0 {
		return
	}
	currentTime := util.GetNow()
	for _, v := range list {
		v.CreateTime = currentTime
		v.UpdateTime = currentTime
	}
	db.Create(list)
}

func (d *AlertContactGroupRelDao) Update(db *gorm.DB, list []*model.AlertContactGroupRel, param form.AlertContactParam) {
	if param.ContactId != "" {
		db.Where("tenant_id = ? AND contact_id = ?", param.TenantId, param.ContactId).Delete(model.AlertContactGroupRel{})
	} else {
		db.Where("tenant_id = ? AND group_id = ?", param.TenantId, param.GroupId).Delete(model.AlertContactGroupRel{})
	}
	d.InsertBatch(db, list)
}

func (d *AlertContactGroupRelDao) Delete(db *gorm.DB, entity *model.AlertContactGroupRel) {
	if entity.ContactId != "" {
		db.Where("tenant_id = ? AND contact_id = ?", entity.TenantId, entity.ContactId).Delete(model.AlertContactGroupRel{})
	} else {
		db.Where("tenant_id = ? AND group_id = ?", entity.TenantId, entity.GroupId).Delete(model.AlertContactGroupRel{})
	}
}

func (d *AlertContactGroupRelDao) GetAlertContact(db *gorm.DB, tenantId string, contactId string, groupId string) *[]form.AlertContactForm {
	var model = &[]form.AlertContactForm{}
	db.Raw(checkAlertContact, groupId, tenantId, contactId).Find(model)
	return model
}

func (d *AlertContactGroupRelDao) GetAlertContactGroup(db *gorm.DB, tenantId string, groupId string) *[]form.AlertContactGroupForm {
	var model = &[]form.AlertContactGroupForm{}
	db.Raw(checkAlertContactGroup, tenantId, groupId).Find(model)
	return model
}
