package dao

import (
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/source_type"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	cvo "code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"fmt"
	"gorm.io/gorm"
	"strings"
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

func (a *AlertRecordDao) GetPageList(db *gorm.DB, tenantId string, f form.AlertRecordPageQueryForm) *cvo.PageVO {
	var list []vo.AlertRecordPageVO
	var c = bytes.Buffer{}
	c.WriteString("SELECT ar.id, ar.`status`, ar.tenant_id, ar.rule_id, ar.rule_name, ar.monitor_type,ar.source_type,ar.source_id, ar.summary,ar.current_value,ar.start_time,ar.end_time,ar.target_value,ar.expression,ar.duration,ar.`level`,ar.notice_status,ar.alarm_key, ar.region,ar.create_time,ar.update_time  FROM t_alert_record ar  WHERE ar.rule_source_type!=? and ar.tenant_id=? ")
	var cv []interface{}
	cv = append(cv, source_type.AutoScaling, tenantId)

	if strutil.IsNotBlank(f.Level) {
		c.WriteString(" and ar.level in (?) ")
		cv = append(cv, strings.Split(f.Level, ","))
	}
	if strutil.IsNotBlank(f.Region) {
		c.WriteString(" and ar.region=? ")
		cv = append(cv, f.Region)
	}
	if strutil.IsNotBlank(f.ResourceId) {
		c.WriteString(" and ar.source_id=? ")
		cv = append(cv, f.ResourceId)
	}
	if strutil.IsNotBlank(f.ResourceType) {
		c.WriteString(" and ar.source_type=? ")
		cv = append(cv, f.ResourceType)
	}
	if strutil.IsNotBlank(f.RuleId) {
		c.WriteString(" and ar.rule_id=? ")
		cv = append(cv, f.RuleId)
	}
	if strutil.IsNotBlank(f.RuleName) {
		c.WriteString(" and ar.rule_name like concat('%', ?, '%') ")
		cv = append(cv, f.RuleName)
	}
	if strutil.IsNotBlank(f.Status) {
		c.WriteString(" and ar.status=? ")
		cv = append(cv, f.Status)
	}
	if strutil.IsNotBlank(f.StartTime) {
		c.WriteString(" and ar.create_time>=? ")
		cv = append(cv, f.StartTime)
	}
	if strutil.IsNotBlank(f.EndTime) {
		c.WriteString(" and ar.create_time<=? ")
		cv = append(cv, f.EndTime)
	}
	var total int64
	db.Raw("select count(1) from ("+c.String()+" ) t", cv...).Scan(&total)
	if total > 0 && int(total) >= (f.PageNum-1)*f.PageSize {
		c.WriteString(" order by ar.create_time desc limit ?,?")
		cv = append(cv, (f.PageNum-1)*f.PageSize)
		cv = append(cv, f.PageSize)
		db.Raw(c.String(), cv...).Find(&list)
	}

	return &cvo.PageVO{
		Records: list,
		Total:   int(total),
		Size:    f.PageSize,
		Current: f.PageNum,
		Pages:   (int(total) / f.PageSize) + 1,
	}
}

func (a *AlertRecordDao) GetByIdAndTenantId(db *gorm.DB, id, tenantId string) *vo.AlertRecordDetailVO {
	var detail vo.AlertRecordDetailVO
	db.Model(commonModels.AlertRecord{}).Where(&commonModels.AlertRecord{Id: id, TenantId: tenantId}).Find(&detail)
	return &detail
}

func (a *AlertRecordDao) GetAlertRecordTotal(db *gorm.DB, tenantId string, region string, startTime string, endTime string) int64 {
	var count int64
	if region != "" {
		db.Model(&commonModels.AlertRecord{}).Where("tenant_id = ? AND status = ? AND create_time BETWEEN ? AND ? AND region = ? and rule_source_type != ?", tenantId, "firing", startTime, endTime, region, source_type.AutoScaling).Count(&count)
	} else {
		db.Model(&commonModels.AlertRecord{}).Where("tenant_id = ? AND status = ? AND create_time BETWEEN ? AND ? and rule_source_type != ?", tenantId, "firing", startTime, endTime, source_type.AutoScaling).Count(&count)
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
