package dao

import (
	"bytes"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	cvo "code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/vo"
	"gorm.io/gorm"
)

type AlertRecordDao struct {
}

var AlertRecord = new(AlertRecordDao)

func (a *AlertRecordDao) GetPageList(db *gorm.DB, tenantId string, f forms.AlertRecordPageQueryForm) *cvo.PageVO {
	var list []vo.AlertRecordPageVO

	var c = bytes.Buffer{}
	c.WriteString("tenant_id=? ")
	var cv []interface{}
	cv = append(cv, tenantId)

	if tools.IsNotBlank(f.Level) {
		c.WriteString(" and level in (?) ")
		cv = append(cv, f.Level)
	}
	if tools.IsNotBlank(f.Region) {
		c.WriteString(" and region=? ")
		cv = append(cv, f.Region)
	}
	if tools.IsNotBlank(f.ResourceId) {
		c.WriteString(" and source_id=? ")
		cv = append(cv, f.ResourceId)
	}
	if tools.IsNotBlank(f.ResourceType) {
		c.WriteString(" and source_type=? ")
		cv = append(cv, f.ResourceType)
	}
	if tools.IsNotBlank(f.RuleId) {
		c.WriteString(" and rule_id=? ")
		cv = append(cv, f.RuleId)
	}
	if tools.IsNotBlank(f.RuleName) {
		c.WriteString(" and rule_name like concat('%', ?, '%') ")
		cv = append(cv, f.RuleName)
	}
	if tools.IsNotBlank(f.Status) {
		c.WriteString(" and status=? ")
		cv = append(cv, f.Status)
	}

	if tools.IsNotBlank(f.StartTime) {
		c.WriteString(" and create_time>=? ")
		cv = append(cv, f.StartTime)
	}
	if tools.IsNotBlank(f.EndTime) {
		c.WriteString(" and create_time<=? ")
		cv = append(cv, f.EndTime)
	}
	var total int64
	db.Model(&commonModels.AlertRecord{}).Where(c.String(), cv...).Count(&total)
	if total > 0 && int(total) >= (f.PageNum-1)*f.PageSize {
		db.Model(&commonModels.AlertRecord{}).Order("create_time desc ").Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Offset((f.PageNum - 1) * f.PageSize).Limit(f.PageSize)
		}).Where(c.String(), cv...).Find(&list)
	}

	return &cvo.PageVO{
		Records: list,
		Total:   int(total),
		Size:    f.PageSize,
		Current: f.PageNum,
		Pages:   (int(total) / f.PageSize) + 1,
	}
}

func (a *AlertRecordDao) GetById(db *gorm.DB, id string) *vo.AlertRecordDetailVO {
	var vo vo.AlertRecordDetailVO
	db.Model(commonModels.AlertRecord{}).Find(&vo, id)
	return &vo
}
