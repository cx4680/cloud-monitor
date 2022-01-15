package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type ContactDao struct {
}

var Contact = new(ContactDao)

const (
	SelectContact = "SELECT " +
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
		"LEFT JOIN " +
		"( " +
		"SELECT " +
		"acgr.contact_biz_id AS contact_biz_id, " +
		"ANY_VALUE( acgr.tenant_id ) AS tenant_id, " +
		"GROUP_CONCAT( acg.biz_id ) AS group_biz_id, " +
		"GROUP_CONCAT( acg.name ) AS group_name, " +
		"COUNT( acgr.contact_biz_id ) AS group_count " +
		"FROM " +
		"t_contact_group AS acg " +
		"LEFT JOIN t_contact_group_rel AS acgr ON acg.biz_id = acgr.group_biz_id AND acg.tenant_id = acgr.tenant_id " +
		"GROUP BY " +
		"acgr.contact_biz_id " +
		") " +
		"AS acg ON ac.biz_id = acg.contact_biz_id AND ac.tenant_id = acg.tenant_id " +
		"WHERE " +
		"ac.state = 1 " +
		"AND ac.tenant_id = ? " +
		"%s" +
		"GROUP BY " +
		"ac.biz_id " +
		"ORDER BY " +
		"ac.create_time DESC "
)

func (d *ContactDao) Select(db *gorm.DB, param form.ContactParam) *form.ContactFormPage {
	var entity = &[]form.ContactForm{}
	var sql string
	var value string
	var total int64
	if strutil.IsNotBlank(param.ContactName) {
		sql = " AND ac.name LIKE CONCAT('%',?,'%') "
		value = param.ContactName
	} else if strutil.IsNotBlank(param.Phone) {
		sql = " AND ac.biz_id = ANY(SELECT contact_biz_id FROM t_contact_information WHERE type = 1 AND state = 1 AND address LIKE CONCAT('%',?,'%')) "
		value = param.Phone
	} else if strutil.IsNotBlank(param.Email) {
		sql = " AND ac.biz_id = ANY(SELECT contact_biz_id FROM t_contact_information WHERE type = 2 AND state = 1 AND address LIKE CONCAT('%',?,'%')) "
		value = param.Email
	}
	sql = fmt.Sprintf(SelectContact, sql)
	db.Raw("select count(1) from ( "+sql+") t ", param.TenantId, value).Scan(&total)
	sql += "LIMIT " + strconv.Itoa((param.PageCurrent-1)*param.PageSize) + "," + strconv.Itoa(param.PageSize)
	db.Raw(sql, param.TenantId, value).Find(entity)
	var contactFormPage = &form.ContactFormPage{
		Records: entity,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return contactFormPage
}

func (d *ContactDao) Insert(db *gorm.DB, entity *model.Contact) {
	db.Create(entity)
}

func (d *ContactDao) Update(db *gorm.DB, entity *model.Contact) {
	db.Model(&model.Contact{}).Where("tenant_id = ? AND biz_id = ?", entity.TenantId, entity.BizId).Updates(entity)
}

func (d *ContactDao) Delete(db *gorm.DB, entity *model.Contact) {
	db.Where("tenant_id = ? AND biz_id = ?", entity.TenantId, entity.BizId).Delete(model.Contact{})
}

func (d *ContactDao) ActivateContact(db *gorm.DB, activeCode string) string {
	var entity = &model.ContactInformation{}
	db.Model(entity).Where("active_code = ?", activeCode).Update("state", 1)
	return entity.TenantId
}
