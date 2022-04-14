package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	model2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"gorm.io/gorm"
	"strconv"
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
		"ac.create_time AS create_time, " +
		"ac.update_time AS update_time, " +
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
	var modelList []form.ContactGroupForm
	var total int64
	db.Raw("select count(1) from ( "+SelectContactGroup+") t ", param.TenantId, param.GroupName).Scan(&total)
	if total >= 0 {
		db.Raw(SelectContactGroup+" LIMIT ?,? ", param.TenantId, param.GroupName, strconv.Itoa((param.PageCurrent-1)*param.PageSize), strconv.Itoa(param.PageSize)).Find(&modelList)
	}
	var contactFormPage = &form.ContactFormPage{
		Records: modelList,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return contactFormPage
}

func (d *ContactGroupDao) SelectGroupContact(db *gorm.DB, param form.ContactParam) *form.ContactFormPage {
	var modelList []form.ContactForm
	var total int64
	db.Model(&model2.ContactGroupRel{}).Where("tenant_id = ? AND group_biz_id = ?", param.TenantId, param.GroupBizId).Count(&total)
	var contactFormPage = &form.ContactFormPage{
		Records: modelList,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	if total == 0 {
		return contactFormPage
	}
	db.Raw(SelectGroupContact+"LIMIT ?,?", param.GroupBizId, param.TenantId, param.GroupBizId, strconv.Itoa((param.PageCurrent-1)*param.PageSize), strconv.Itoa(param.PageSize)).Find(&modelList)
	contactFormPage.Records = modelList
	return contactFormPage
}

func (d *ContactGroupDao) Insert(db *gorm.DB, entity *model2.ContactGroup) {
	db.Create(entity)
}

func (d *ContactGroupDao) Update(db *gorm.DB, entity *model2.ContactGroup) {
	db.Save(entity)
}

func (d *ContactGroupDao) Delete(db *gorm.DB, entity *model2.ContactGroup) {
	db.Model(&model2.ContactGroup{}).Where("tenant_id = ? AND biz_id = ?", entity.TenantId, entity.BizId).Delete(model2.ContactGroup{})
}

//查询组
func (d *ContactGroupDao) GetGroup(tenantId, groupBizId string) model2.ContactGroup {
	var contactGroup model2.ContactGroup
	global.DB.Where("tenant_id = ? AND biz_id = ?", tenantId, groupBizId).First(&contactGroup)
	return contactGroup
}

//查询租户下的联系组数量
func (d *ContactGroupDao) GetGroupCount(tenantId string) int64 {
	var groupCount int64
	global.DB.Model(&model2.ContactGroup{}).Where("tenant_id = ?", tenantId).Count(&groupCount)
	return groupCount
}

//校验组名是否重复
func (d *ContactGroupDao) CheckGroupName(tenantId, groupName, groupBizId string) bool {
	var count int64
	if strutil.IsBlank(groupBizId) {
		global.DB.Model(&model2.ContactGroup{}).Where("tenant_id = ? AND name = ?", tenantId, groupName).Count(&count)
	} else {
		global.DB.Model(&model2.ContactGroup{}).Where("tenant_id = ? AND name = ? AND biz_id != ?", tenantId, groupName, groupBizId).Count(&count)
	}
	if count > 0 {
		return true
	}
	return false
}

//检验组ID是否存在
func (d *ContactGroupDao) CheckGroupId(tenantId, groupBizId string) bool {
	var count int64
	global.DB.Model(&model2.ContactGroup{}).Where("tenant_id = ? AND biz_id = ?", tenantId, groupBizId).Count(&count)
	if count > 0 {
		return true
	}
	return false
}
