package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"gorm.io/gorm"
)

type ContactGroupDao struct {
}

var ContactGroup = new(ContactGroupDao)

const (
	SelectContactGroup = "SELECT " +
		"acg.biz_id AS group_biz_id, " +
		"acg.name AS group_name, " +
		"acg.description AS description, " +
		"acg.create_time AS create_time, " +
		"acg.update_time AS update_time, " +
		"COUNT( acgr.group_biz_id ) AS contact_count " +
		"FROM " +
		"t_contact_group AS acg " +
		"LEFT JOIN t_contact_group_rel AS acgr ON acg.biz_id = acgr.group_biz_id AND acg.tenant_id = acgr.tenant_id " +
		"WHERE " +
		"acg.tenant_id = ? " +
		"AND acg.name LIKE CONCAT('%',?,'%') " +
		"GROUP BY " +
		"acg.biz_id " +
		"ORDER BY " +
		"acg.create_time DESC "

	SelectGroupContact = "SELECT " +
		"ac.biz_id AS contact_biz_id, " +
		"ac.name AS contact_name, " +
		"ANY_VALUE( acg.group_biz_id ) AS group_biz_id, " +
		"ANY_VALUE( acg.group_name ) AS group_name, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.address END ) AS phone, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.state END ) AS phone_state, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.address END ) AS email, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.state END ) AS email_state, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.address END ) AS lanxin, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.state END ) AS lanxin_state, " +
		"ac.description AS description, " +
		"ANY_VALUE( acg.group_count ) AS group_count " +
		"FROM " +
		"t_contact AS ac " +
		"LEFT JOIN t_contact_information AS aci ON ac.biz_id = aci.contact_biz_id AND ac.tenant_id = aci.tenant_id " +
		"LEFT JOIN ( " +
		"SELECT " +
		"acgr.contact_biz_id AS contact_biz_id, " +
		"ANY_VALUE( acgr.tenant_id ) AS tenant_id, " +
		"GROUP_CONCAT( acg.biz_id ) AS group_biz_id, " +
		"GROUP_CONCAT( acg.name ) AS group_name, " +
		"COUNT( acgr.contact_biz_id ) AS group_count " +
		"FROM " +
		"t_contact_group AS acg " +
		"LEFT JOIN t_contact_group_rel AS acgr ON acg.biz_id = acgr.group_biz_id AND acg.tenant_id = acgr.tenant_id " +
		"WHERE acg.biz_id = ? " +
		"GROUP BY " +
		"acgr.contact_biz_id ) " +
		"AS acg ON ac.biz_id = acg.contact_biz_id AND ac.tenant_id = acg.tenant_id " +
		"WHERE " +
		"ac.state = 1 " +
		"AND ac.tenant_id = ? " +
		"AND acg.group_biz_id = ? " +
		"GROUP BY " +
		"ac.biz_id " +
		"ORDER BY " +
		"ac.create_time DESC "
)

func (d *ContactGroupDao) SelectContactGroup(db *gorm.DB, param form.ContactParam) *form.ContactFormPage {
	var modelList = &[]form.ContactGroupForm{}
	var total int64
	db.Raw("select count(1) from ( "+SelectContactGroup+") t ", param.TenantId, param.GroupName).Scan(&total)
	db.Raw(SelectContactGroup, param.TenantId, param.GroupName).Find(modelList)
	var contactFormPage = &form.ContactFormPage{
		Records: modelList,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return contactFormPage
}

func (d *ContactGroupDao) SelectGroupContact(db *gorm.DB, param form.ContactParam) *form.ContactFormPage {
	var modelList = &[]form.ContactForm{}
	var total int64
	db.Model(&model.ContactGroupRel{}).Where("tenant_id = ? AND group_biz_id = ?", param.TenantId, param.GroupBizId).Count(&total)
	var contactFormPage = &form.ContactFormPage{
		Records: modelList,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	if total == 0 {
		return contactFormPage
	}
	db.Raw(SelectGroupContact, param.GroupBizId, param.TenantId, param.GroupBizId).Find(modelList)
	contactFormPage.Records = modelList
	return contactFormPage
}

func (d *ContactGroupDao) Insert(db *gorm.DB, entity *model.ContactGroup) {
	db.Create(entity)
}

func (d *ContactGroupDao) Update(db *gorm.DB, entity *model.ContactGroup) {
	db.Save(entity)
}

func (d *ContactGroupDao) Delete(db *gorm.DB, entity *model.ContactGroup) {
	db.Model(&model.ContactGroup{}).Where("tenant_id = ? AND biz_id = ?", entity.TenantId, entity.BizId).Delete(model.ContactGroup{})
}
