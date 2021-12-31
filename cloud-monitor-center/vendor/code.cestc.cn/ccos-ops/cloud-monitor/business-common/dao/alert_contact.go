package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactDao struct {
}

var AlertContact = new(AlertContactDao)

const (
	SelectAlterContact = "SELECT " +
		"ac.id AS contact_id, " +
		"ac.name AS contact_name, " +
		"ANY_VALUE( acg.group_id ) AS group_id, " +
		"ANY_VALUE( acg.group_name ) AS group_name, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.NO END ) AS phone, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.is_certify END ) AS phone_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.NO END ) AS email, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.is_certify END ) AS email_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.NO END ) AS lanxin, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.is_certify END ) AS lanxin_certify, " +
		"ac.description AS description, " +
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
		"COUNT( acgr.contact_id ) AS group_count " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id AND acg.tenant_id = acgr.tenant_id " +
		"GROUP BY " +
		"acgr.contact_id " +
		") " +
		"AS acg ON ac.id = acg.contact_id AND ac.tenant_id = acg.tenant_id " +
		"WHERE " +
		"ac.status = 1 " +
		"AND ac.tenant_id = %s " +
		"%s" +
		"GROUP BY " +
		"ac.id " +
		"ORDER BY " +
		"ac.create_time DESC "
)

func (d *AlertContactDao) Select(db *gorm.DB, param form.AlertContactParam) *form.AlertContactFormPage {
	var model = &[]form.AlertContactForm{}
	var sql string
	if strutil.IsNotBlank(param.ContactName) {
		sql += " AND ac.name LIKE CONCAT('%','" + param.ContactName + "','%') "
	}
	if strutil.IsNotBlank(param.Phone) {
		sql += " AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 1 AND is_certify = 1 AND no LIKE CONCAT('%','" + param.Phone + "','%')) "
	}
	if strutil.IsNotBlank(param.Email) {
		sql += " AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 2 AND is_certify = 1 AND no LIKE CONCAT('%','" + param.Email + "','%')) "
	}
	sql = fmt.Sprintf(SelectAlterContact, param.TenantId, sql)
	var total int64
	db.Raw("select count(1) from ( " + sql + ") t ").Scan(&total)
	sql += "LIMIT " + strconv.Itoa((param.PageCurrent-1)*param.PageSize) + "," + strconv.Itoa(param.PageSize)
	db.Raw(sql).Find(model)
	var alertContactFormPage = &form.AlertContactFormPage{
		Records: model,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return alertContactFormPage
}

func (d *AlertContactDao) Insert(db *gorm.DB, entity *model.AlertContact) {
	db.Create(entity)
}

func (d *AlertContactDao) Update(db *gorm.DB, entity *model.AlertContact) {
	db.Updates(entity)
}

func (d *AlertContactDao) Delete(db *gorm.DB, entity *model.AlertContact) {
	db.Where("tenant_id = ? AND id = ?", entity.TenantId, entity.Id).Delete(model.AlertContact{})
}

func (d *AlertContactDao) CertifyAlertContact(db *gorm.DB, activeCode string) string {
	var model = &model.AlertContactInformation{}
	db.Model(model).Where("active_code = ?", activeCode).Update("is_certify", 1)
	return model.TenantId
}
