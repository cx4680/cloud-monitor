package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
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
		"ANY_VALUE( acg.group_count ) AS group_count, " +
		"ac.create_time AS create_time, " +
		"ac.update_time AS update_time " +
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
	var entity []form.ContactForm
	var sql string
	var screenParams []interface{}
	var total int64
	screenParams = append(screenParams, param.TenantId)
	if strutil.IsNotBlank(param.ContactName) {
		sql += " AND ac.name LIKE CONCAT('%',?,'%') "
		screenParams = append(screenParams, param.ContactName)
	}
	if strutil.IsNotBlank(param.Phone) {
		sql += " AND ac.biz_id = ANY(SELECT contact_biz_id FROM t_contact_information WHERE type = 1 AND state = 1 AND address LIKE CONCAT('%',?,'%')) "
		screenParams = append(screenParams, param.Phone)
	}
	if strutil.IsNotBlank(param.Email) {
		sql += " AND ac.biz_id = ANY(SELECT contact_biz_id FROM t_contact_information WHERE type = 2 AND state = 1 AND address LIKE CONCAT('%',?,'%')) "
		screenParams = append(screenParams, param.Email)
	}
	sql = fmt.Sprintf(SelectContact, sql)
	db.Raw(fmt.Sprintf("SELECT COUNT(1) FROM (%s) t ", sql), screenParams...).Scan(&total)
	if total >= 0 {
		sql = fmt.Sprintf("%s LIMIT %s,%s", sql, strconv.Itoa((param.PageCurrent-1)*param.PageSize), strconv.Itoa(param.PageSize))
		db.Raw(sql, screenParams...).Find(&entity)
	}
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
	db.Where("tenant_id = ? AND biz_id = ? AND state = 1", entity.TenantId, entity.BizId).Updates(entity)
}

func (d *ContactDao) Delete(db *gorm.DB, entity *model.Contact) {
	db.Select("state", "update_time").Where("tenant_id = ? AND biz_id = ?", entity.TenantId, entity.BizId).Updates(entity)
}

func (d *ContactDao) ActivateContact(db *gorm.DB, activeCode string) string {
	var entity = &model.ContactInformation{
		UpdateTime: util.GetNow(),
		State:      1,
	}
	db.Select("state", "update_time").Where("active_code = ?", activeCode).Updates(entity).Find(entity)
	return entity.TenantId
}

// GetContactCount 查询租户下的联系人数量
func (d *ContactDao) GetContactCount(tenantId string) int64 {
	var count int64
	global.DB.Model(&model.Contact{}).Where("tenant_id = ? AND state = 1", tenantId).Count(&count)
	return count
}

// CheckContact 校验联系人
func (d *ContactDao) CheckContact(tenantId, contactBizId string) bool {
	var count int64
	global.DB.Model(&model.Contact{}).Where("tenant_id = ? AND biz_id = ? AND state = 1", tenantId, contactBizId).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

// GetTenantIdByActiveCode 根据激活码获取租户ID
func (d *ContactDao) GetTenantIdByActiveCode(activeCode string) string {
	var entity = &model.ContactInformation{}
	global.DB.Where("active_code = ?", activeCode).First(entity)
	return entity.TenantId
}

func (d *ContactDao) GetContactBizIdByName(tenantId, contactName string) string {
	var contact = &model.Contact{}
	global.DB.Where("tenant_id = ? AND name = ? AND state = 1", tenantId, contactName).First(contact)
	if contact == (&model.Contact{}) {
		return ""
	}
	return contact.BizId
}
