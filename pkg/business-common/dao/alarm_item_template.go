package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"gorm.io/gorm"
)

type AlarmItemTemplateDao struct{}

var AlarmItemTemplate = new(AlarmItemTemplateDao)

func (dao *AlarmItemTemplateDao) QueryItemListByTemplate(db *gorm.DB, templateBizId string) []model.AlarmItemTemplate {
	var list []model.AlarmItemTemplate
	db.Where(&model.AlarmItemTemplate{RuleTemplateBizId: templateBizId}).Find(&list)
	return list
}
