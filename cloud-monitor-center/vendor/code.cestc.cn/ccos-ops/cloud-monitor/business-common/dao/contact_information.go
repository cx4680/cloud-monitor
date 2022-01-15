package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"gorm.io/gorm"
)

type ContactInformationDao struct {
}

var ContactInformation = new(ContactInformationDao)

func (d *ContactInformationDao) Insert(db *gorm.DB, entity *model.ContactInformation) {
	currentTime := util.GetNow()
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	db.Create(entity)
}

func (d *ContactInformationDao) InsertBatch(db *gorm.DB, list []*model.ContactInformation) {
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

func (d *ContactInformationDao) Update(db *gorm.DB, list []*model.ContactInformation) {
	if len(list) == 0 {
		return
	}
	var newList []*model.ContactInformation
	for _, v := range list {
		db.Where("tenant_id = ? AND contact_biz_id = ? AND type = ?", v.TenantId, v.ContactBizId, v.Type).Delete(&model.ContactInformation{})
		if strutil.IsNotBlank(v.Address) {
			newList = append(newList, v)
		}
	}
	d.InsertBatch(db, newList)
}

func (d *ContactInformationDao) Delete(db *gorm.DB, entity *model.ContactInformation) {
	db.Model(&model.ContactInformation{}).Where("tenant_id = ? AND contact_biz_id = ?", entity.TenantId, entity.ContactBizId).Delete(&model.ContactInformation{})
}
