package dao

import (
	"bytes"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	cvo "code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/vo"
	"fmt"
	"gorm.io/gorm"
)

type AlertRecordDao struct {
}

var AlertRecord = new(AlertRecordDao)

const (
	recordNumHistory = "SELECT COUNT(t.id) AS number, " +
		"DATE_FORMAT(t.create_time, '%s') AS DayTime " +
		"FROM t_alert_record t " +
		"WHERE t.status = 'firing' " +
		"AND t.tenant_id = %s " +
		"AND t.create_time between '%s' AND '%s' " +
		"%s " +
		"GROUP BY daytime"
)

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

func (a *AlertRecordDao) GetAlertRecordTotal(db *gorm.DB, tenantId string, region string, startTime string, endTime string) int64 {
	var count int64
	if region != "" {
		db.Model(&commonModels.AlertRecord{}).Where("tenant_id = ? AND status = ? AND create_time BETWEEN ? AND ? AND region = ?", tenantId, "firing", startTime, endTime, region).Count(&count)
	} else {
		db.Model(&commonModels.AlertRecord{}).Where("tenant_id = ? AND status = ? AND create_time BETWEEN ? AND ?", tenantId, "firing", startTime, endTime).Count(&count)
	}
	return count
}

func (a *AlertRecordDao) GetRecordNumHistory(db *gorm.DB, tenantId string, region string, startTime string, endTime string) []vo.RecordNumHistory {
	var sql string
	if region != "" {
		sql = fmt.Sprintf(recordNumHistory, "%Y-%m-%d", tenantId, startTime, endTime, " AND t.region = "+region)
	} else {
		sql = fmt.Sprintf(recordNumHistory, "%Y-%m-%d", tenantId, startTime, endTime, "")
	}
	var recordNumHistory []vo.RecordNumHistory
	db.Raw(sql).Find(&recordNumHistory)
	return recordNumHistory
}
