package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"gorm.io/gorm"
)

type AlertContactGroupDao struct {
}

var AlertContactGroup = new(AlertContactGroupDao)

const (
	SelectAlterContactGroup = "SELECT " +
		"acg.id AS group_id, " +
		"acg.name AS group_name, " +
		"acg.description AS description, " +
		"acg.create_time AS create_time, " +
		"acg.update_time AS update_time, " +
		"COUNT( acgr.group_id ) AS contact_count " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id AND acg.tenant_id = acgr.tenant_id " +
		"WHERE " +
		"acg.tenant_id = ? " +
		"AND acg.name LIKE CONCAT('%',?,'%') " +
		"GROUP BY " +
		"acg.id " +
		"ORDER BY " +
		"acg.create_time DESC "

	SelectAlterGroupContact = "SELECT " +
		"ac.id AS contact_id, " +
		"ac.name AS contact_name, " +
		"acg.group_id AS group_id, " +
		"acg.group_name AS group_name, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.NO END ) AS phone, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.is_certify END ) AS phone_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.NO END ) AS email, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.is_certify END ) AS email_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.NO END ) AS lanxin, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.is_certify END ) AS lanxin_certify, " +
		"ac.description AS description " +
		"FROM " +
		"alert_contact AS ac " +
		"LEFT JOIN alert_contact_information AS aci ON ac.id = aci.contact_id AND ac.tenant_id = aci.tenant_id " +
		"LEFT JOIN ( " +
		"SELECT " +
		"acgr.contact_id AS contact_id, " +
		"acgr.tenant_id AS tenant_id, " +
		"GROUP_CONCAT( acg.id ) AS group_id, " +
		"GROUP_CONCAT( acg.name ) AS group_name " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id AND acg.tenant_id = acgr.tenant_id " +
		"GROUP BY " +
		"acgr.contact_id ) " +
		"AS acg ON ac.id = acg.contact_id AND ac.tenant_id = acg.tenant_id " +
		"WHERE " +
		"ac.status = 1 " +
		"AND ac.tenant_id = ? " +
		"AND acg.group_id = ? " +
		"GROUP BY " +
		"ac.id " +
		"ORDER BY " +
		"ac.create_time DESC "
)

func (d *AlertContactGroupDao) SelectAlertContactGroup(db *gorm.DB, param forms.AlertContactParam) *forms.AlertContactFormPage {
	var model = &[]forms.AlertContactGroupForm{}
	var total int64
	db.Model(&models.AlertContactGroup{}).Where("tenant_id = ?", param.TenantId).Count(&total)
	db.Raw(SelectAlterContactGroup, param.TenantId, param.GroupName).Find(model)
	var alertContactFormPage = &forms.AlertContactFormPage{
		Records: model,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return alertContactFormPage
}

func (d *AlertContactGroupDao) SelectAlertGroupContact(db *gorm.DB, param forms.AlertContactParam) *[]forms.AlertContactForm {
	var model = &[]forms.AlertContactForm{}
	db.Raw(SelectAlterGroupContact, param.TenantId, param.GroupId).Find(model)
	return model
}

func (d *AlertContactGroupDao) Insert(db *gorm.DB, entity *models.AlertContactGroup) {
	db.Create(entity)
}

func (d *AlertContactGroupDao) Update(db *gorm.DB, entity *models.AlertContactGroup) {
	db.Updates(entity)
}

func (d *AlertContactGroupDao) Delete(db *gorm.DB, entity *models.AlertContactGroup) {
	db.Where("tenant_id = ? AND id = ?", entity.TenantId, entity.Id).Delete(models.AlertContactGroup{})
}
