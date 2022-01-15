package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"gorm.io/gorm"
)

type ContactGroupRelDao struct {
}

const (
	checkContact = "SELECT " +
		"ac.biz_id AS contact_biz_id, " +
		"ac.name AS contact_name, " +
		"ANY_VALUE( acg.group_biz_id ) AS group_biz_id, " +
		"ANY_VALUE( acg.group_name ) AS group_name, " +
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
		"COUNT( acgr.contact_biz_id ) AS group_count  " +
		"FROM " +
		"t_contact_group AS acg " +
		"LEFT JOIN t_contact_group_rel AS acgr ON acg.biz_id = acgr.group_biz_id AND acg.tenant_id = acgr.tenant_id " +
		"WHERE acg.biz_id != ? " +
		"GROUP BY " +
		"acgr.contact_biz_id " +
		") " +
		"AS acg ON ac.biz_id = acg.contact_biz_id AND ac.tenant_id = acg.tenant_id " +
		"WHERE " +
		"ac.state = 1 " +
		"AND ac.tenant_id = ? " +
		"AND ac.biz_id = ? " +
		"GROUP BY " +
		"ac.biz_id "

	checkContactGroup = "SELECT " +
		"acg.biz_id AS group_biz_id, " +
		"acg.name AS group_name, " +
		"COUNT( acgr.group_biz_id ) AS contact_count " +
		"FROM " +
		"t_contact_group AS acg " +
		"LEFT JOIN t_contact_group_rel AS acgr ON acg.biz_id = acgr.group_biz_id AND acg.tenant_id = acgr.tenant_id " +
		"WHERE " +
		"acg.tenant_id = ? " +
		"AND acg.biz_id = ? " +
		"GROUP BY " +
		"acg.biz_id "
)

var ContactGroupRel = new(ContactGroupRelDao)

func (d *ContactGroupRelDao) Insert(db *gorm.DB, entity *model.ContactGroupRel) {
	currentTime := util.GetNow()
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (d *ContactGroupRelDao) InsertBatch(db *gorm.DB, list []*model.ContactGroupRel) {
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

func (d *ContactGroupRelDao) Update(db *gorm.DB, list []*model.ContactGroupRel, param form.ContactParam) {
	if strutil.IsNotBlank(param.ContactBizId) {
		db.Where("tenant_id = ? AND contact_biz_id = ?", param.TenantId, param.ContactBizId).Delete(model.ContactGroupRel{})
	} else {
		db.Where("tenant_id = ? AND group_biz_id = ?", param.TenantId, param.GroupBizId).Delete(model.ContactGroupRel{})
	}
	d.InsertBatch(db, list)
}

func (d *ContactGroupRelDao) Delete(db *gorm.DB, entity *model.ContactGroupRel) {
	if strutil.IsNotBlank(entity.ContactBizId) {
		db.Model(&model.ContactGroupRel{}).Where("tenant_id = ? AND contact_biz_id = ?", entity.TenantId, entity.ContactBizId).Delete(model.ContactGroupRel{})
	} else {
		db.Model(&model.ContactGroupRel{}).Where("tenant_id = ? AND group_biz_id = ?", entity.TenantId, entity.GroupBizId).Delete(model.ContactGroupRel{})
	}
}

func (d *ContactGroupRelDao) GetContact(db *gorm.DB, tenantId string, contactBizId string, groupBizId string) *[]form.ContactForm {
	var entity = &[]form.ContactForm{}
	db.Raw(checkContact, groupBizId, tenantId, contactBizId).Find(entity)
	return entity
}

func (d *ContactGroupRelDao) GetContactGroup(db *gorm.DB, tenantId string, groupBizId string) *[]form.ContactGroupForm {
	var entity = &[]form.ContactGroupForm{}
	db.Raw(checkContactGroup, tenantId, groupBizId).Find(entity)
	return entity
}
