package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"gorm.io/gorm"
)

type TenantAlarmTemplateRelDao struct {
}

var TenantAlarmTemplateRel = new(TenantAlarmTemplateRelDao)

func (dao *TenantAlarmTemplateRelDao) Insert(db *gorm.DB, record model.TenantAlarmTemplateRel) {
	db.Create(&record)
}

func (dao *TenantAlarmTemplateRelDao) Delete(db *gorm.DB, tenantId string, templateBizIds []string) {
	db.Where("tenant_id=? and template_biz_id in (?)", tenantId, templateBizIds).Delete(&model.TenantAlarmTemplateRel{})
}
